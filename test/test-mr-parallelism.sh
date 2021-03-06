#########################################################

# comment this out to run the tests without the Go race detector.
RACE=-race

if [[ "$OSTYPE" = "darwin"* ]]
then
  if go version | grep 'go1.17.[012345]'
  then
    # -race with plug-ins on x86 MacOS 12 with
    # go1.17 before 1.17.6 sometimes crash.
    RACE=
    echo '*** Turning off -race since it may not work on a Mac'
    echo '    with ' `go version`
  fi
fi

TIMEOUT=timeout
if timeout 2s sleep 1 > /dev/null 2>&1
then
  :
else
  if gtimeout 2s sleep 1 > /dev/null 2>&1
  then
    TIMEOUT=gtimeout
  else
    # no timeout command
    TIMEOUT=
    echo '*** Cannot find timeout command; proceeding without timeouts.'
  fi
fi
if [ "$TIMEOUT" != "" ]
then
  TIMEOUT+=" -k 2s 180s "
fi

rm -rf ../mr-tmp
mkdir ../mr-tmp || exit 1
cd ../mr-tmp || exit 1
rm -f mr-*

(cd ../plugins && go clean)
(cd ../cmd && go clean)
(cd ../plugins/rtiming && go build $RACE -buildmode=plugin rtiming.go) || exit 1
(cd ../plugins/mtiming && go build $RACE -buildmode=plugin mtiming.go) || exit 1
(cd ../cmd/coordinator && go build $RACE mrcoordinator.go) || exit 1
(cd ../cmd/worker && go build $RACE mrworker.go) || exit 1
(cd ../cmd/sequential && go build $RACE mrsequential.go) || exit 1

failed_any=0

cd ../mr-tmp

#########################################################
echo '***' Starting map parallelism test.

rm -f mr-*

$TIMEOUT ../cmd/coordinator/mrcoordinator ../test/testdata/pg*txt &
sleep 1

$TIMEOUT ../cmd/worker/mrworker ../plugins/mtiming/mtiming.so &
$TIMEOUT ../cmd/worker/mrworker ../plugins/mtiming/mtiming.so

NT=`cat mr-out* | grep '^times-' | wc -l | sed 's/ //g'`
if [ "$NT" != "2" ]
then
  echo '---' saw "$NT" workers rather than 2
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

if cat mr-out* | grep '^parallel.* 2' > /dev/null
then
  echo '---' map parallelism test: PASS
else
  echo '---' map workers did not run in parallel
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

wait


#########################################################
echo '***' Starting reduce parallelism test.

rm -f mr-*

$TIMEOUT ../cmd/coordinator/mrcoordinator ../test/testdata/pg*txt &
sleep 1

$TIMEOUT ../cmd/worker/mrworker ../plugins/rtiming/rtiming.so &
$TIMEOUT ../cmd/worker/mrworker ../plugins/rtiming/rtiming.so

NT=`cat mr-out* | grep '^[a-z] 2' | wc -l | sed 's/ //g'`
if [ "$NT" -lt "2" ]
then
  echo '---' saw "$NT" workers rather than 2
  echo '---' too few parallel reduces.
  echo '---' reduce parallelism test: FAIL
  failed_any=1
else
  echo '---' reduce parallelism test: PASS
fi

wait