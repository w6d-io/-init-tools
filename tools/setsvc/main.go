package main

import (
    "context"
    "fmt"
    "os"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    k8s "k8s.io/client-go/kubernetes"

    "github.com/cakturk/go-netstat/netstat"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/klog/v2/klogr"
)

const (
    protoIPv4 = 0x01
    protoIPv6 = 0x02
)

var (
    namespace string
    svcName   string
    log       = klogr.New()
)

func main() {
    if len(os.Args) != 2 {
        usage()
    }
    svcName = os.Args[1]
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
        if err := setsvc(svcName, port); err != nil {
            log.Error(err, "set kubernetes service")
            os.Exit(1)
        }
        os.Exit(0)
    }
    log.Error(fmt.Errorf("setsvc"), "Port not found")
}

func usage() {
    fmt.Println("setsvc : Usage")
    fmt.Println("\tsetsvc <service name>")
    os.Exit(1)
}

func displaySockInfo(_ string, s []netstat.SockTabEntry) uint16 {
    lookup := func(skAddr *netstat.SockAddr) uint16 {
        const IPv4StrLen = 17
        addr := skAddr.IP.String()

        if len(addr) > IPv4StrLen {
            addr = addr[:IPv4StrLen]
        }
        return skAddr.Port
    }
    for _, e := range s {
        return lookup(e.LocalAddr)
    }
    return 0
}

func newClient() (k8s.Interface, error) {
    config, err := getConfig()
    if err != nil {
        log.Error(err, "failed to get k8s config")
        return nil, err
    }
    k8scs, err := k8s.NewForConfig(config)
    if err != nil {
        log.Error(err, "failed to create k8s client from config")
        return nil, err
    }
    return k8scs, nil
}

func getConfig() (config *rest.Config, err error) {
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    configOverrides := &clientcmd.ConfigOverrides{}
    kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
    namespace, _, err = kubeConfig.Namespace()
    if err != nil {
        log.Error(err, "Couldn't get kubeConfiguration namespace")
    }
    config, err = kubeConfig.ClientConfig()
    if err != nil {
        log.Error(err, "Parsing kubeconfig failed")
    }
    return config, nil
}

func setsvc(svcName string, port uint16) (err error) {
    ctx := context.Background()
    client, err := newClient()
    if err != nil {
        log.Error(err, "new client return %v")
    }
    retry := 5
    for {
        svc, err := client.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
        if err != nil {
            log.Error(err, "get svc return")
        }

        log.V(2).Info(fmt.Sprintf("target contains %v", svc.Spec.Ports[0].TargetPort.IntVal))
        svc.Spec.Ports[0].TargetPort.StrVal = ""
        svc.Spec.Ports[0].TargetPort.IntVal = int32(port)
        log.V(2).Info(fmt.Sprintf("target set by %v", svc.Spec.Ports[0].TargetPort.IntVal))
        _, err = client.CoreV1().Services(namespace).Update(ctx, svc, metav1.UpdateOptions{})
        if err != nil {
            log.Error(err, "update svc return")
        }
        rsp, _ := client.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
        if rsp.Spec.Ports[0].TargetPort.IntVal == svc.Spec.Ports[0].TargetPort.IntVal {
            log.Info("port has been set")
            break
        }
        if retry <= 0 {
            log.Error(err, "port update failed")
        }
        log.Info("port has not been set")
        log.V(1).Info("retrying...")
        time.Sleep(2 * time.Second)
        retry--
    }
    return nil
}
