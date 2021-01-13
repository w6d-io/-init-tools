package main

import (
    "fmt"
    "os"
    "time"

    "github.com/cakturk/go-netstat/netstat"
    "gitlab.w6d.io/w6d/library/log"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    k8s "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

const (
    protoIPv4 = 0x01
    protoIPv6 = 0x02
)

var (
    namespace string
    svcname   string
)

func init() {
    log.SetLevel(log.TRACE)
}

func main() {
    if len(os.Args) != 2 {
        usage()
    }
    svcname = os.Args[1]
    var proto uint
    proto |= protoIPv4
    proto |= protoIPv6
    var fn netstat.AcceptFn
    fn = func(*netstat.SockTabEntry) bool { return true }

    var (
    	port uint16
    	retry = 60
    	)
    for {
		tabs, err := netstat.TCPSocks(fn)
        if err == nil {
            port = displaySockInfo("tcp", tabs)
        }
        if port == 0 {
            tabs, err = netstat.TCP6Socks(fn)
            if err == nil {
                port = displaySockInfo("tcp6", tabs)
            }
        }
        if port != 0 && retry > 0 {
			break
        }
		retry--
        time.Sleep(10 * time.Second)
    }
    if port != 0 {
        setsvc(svcname, port)
        os.Exit(0)
    }
    log.Fatal("Port not found")
}

func usage() {
    fmt.Println("setsvc : Usage")
    fmt.Println("\tsetsvc <service name>")
    os.Exit(1)
}

func displaySockInfo(proto string, s []netstat.SockTabEntry) uint16 {
    lookup := func(skaddr *netstat.SockAddr) uint16 {
        const IPv4Strlen = 17
        addr := skaddr.IP.String()

        if len(addr) > IPv4Strlen {
            addr = addr[:IPv4Strlen]
        }
        return skaddr.Port
    }
    for _, e := range s {
        return lookup(e.LocalAddr)
    }
    return 0
}

func newClient() (k8s.Interface, error) {
    config, err := getConfig()
    if err != nil {
        log.Fatal("failed to get k8s config")
    }
    k8scs, err := k8s.NewForConfig(config)
    if err != nil {
        log.Fatal("failed to create k8s client from config")
    }
    return k8scs, nil
}

func getConfig() (config *rest.Config, err error) {
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    configOverrides := &clientcmd.ConfigOverrides{}
    kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
    namespace, _, err = kubeConfig.Namespace()
    if err != nil {
        log.Fatal("Couldn't get kubeConfiguration namespace")
    }
    config, err = kubeConfig.ClientConfig()
    if err != nil {
        log.Fatal("Parsing kubeconfig failed")
    }
    return config, nil
}

func setsvc(svcname string, port uint16) (err error) {

    client, err := newClient()
    if err != nil {
        log.Fatalf("new client return %v", err)
    }
    retry := 5
    for {
        svc, err := client.CoreV1().Services(namespace).Get(svcname, metav1.GetOptions{})
        if err != nil {
            log.Fatalf("get svc return %v", err)
        }
        log.Tracef("target contains %v", svc.Spec.Ports[0].TargetPort.IntVal)
        svc.Spec.Ports[0].TargetPort.StrVal = ""
        svc.Spec.Ports[0].TargetPort.IntVal = int32(port)
        log.Tracef("target set by %v", svc.Spec.Ports[0].TargetPort.IntVal)
        _, err = client.CoreV1().Services(namespace).Update(svc)
        if err != nil {
            log.Fatalf("update svc return %v", err)
        }
        rsp, _ := client.CoreV1().Services(namespace).Get(svcname, metav1.GetOptions{})
        if rsp.Spec.Ports[0].TargetPort.IntVal == svc.Spec.Ports[0].TargetPort.IntVal {
            log.Info("port has been set")
            break
        }
        if retry <= 0 {
            log.Fatal("port update failed")
        }
        log.Warning("port has not been set")
        log.Debug("retrying...")
        time.Sleep(2 * time.Second)
        retry--
    }
    return nil
}
