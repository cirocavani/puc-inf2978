export GOPATH=`pwd`
#export GOROOT=/home/cavani/Software/go-1.3
#export PATH=/home/cavani/Software/go-1.3/bin:$PATH

export CGO_CFLAGS="-I ${GUROBI_HOME}/include"
export CGO_LDFLAGS="-L ${GUROBI_HOME}/lib -lgurobi56"
