Ajustar o arquivo env.sh para o ambiente local.

    source env.sh
    ln -s <diretório_com_problemas> data
    go install parallax/tool/player
    ./bin/player

Parâmetros:

    ./bin/player -help

Servidor:

https://github.com/ExpLog/game-theory-master

...

Outras ferramentas:

    (Calcula fluxo usando Gurobi em uma determinada Instância)
    go install parallax/tool/gurobi
    ./bin/gurobi -help

    (Executa Engine de Bid em uma determinada Instância)
    go install parallax/tool/engine
    ./bin/engine -help
