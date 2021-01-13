- Resolution des problemes d'affichage fonction/fichier qui ne donnaient pas toujours la fonction dans laquelle le log etait appele
- Ajout d'un niveau supplementaire de log
``log.Trace()``
- Ajout de fonctions pour changer les couleurs des logs
 ``log.SetInfoColor(log.Red)``
 ``log.SetDebugColor(log.Reset)``
`` ...``

- Ajout d'une fonction permettant de recuperer et de parser les variables d'environnement (OUTPUT_FORMAT et SET_LEVEL_LOG)
``log.SetConfig()``
- Ajout d'une fonction fonction permettant de recuperer le path complet du fichier
``log.SetPathFunc(true)``
- Les logs dans les go routines anynome affichent desormais le pere
``function=Consume.func1``