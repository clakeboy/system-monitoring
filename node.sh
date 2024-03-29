#! /bin/sh
#程序启动脚本
### BEGIN INIT INFO
# Provides:          sys-monitor
# Required-Start:    $remote_fs $network
# Required-Stop:     $remote_fs $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts sys-monitor
# Description:       sys-monitor
### END INIT INFO
prefix=`pwd`
exec_prefix=${prefix}/sys-monitor
pid_file=${prefix}/sys-monitor.pid
conf_file=${prefix}/main.conf

wait_for_pid () {
    try=0

    while test $try -lt 35 ; do

        case "$1" in
            'created')
            if [ -f "$2" ] ; then
                try=''
                break
            fi
            ;;

            'removed')
            if [ ! -f "$2" ] ; then
                try=''
                break
            fi
            ;;
        esac

        echo -n .
        try=`expr $try + 1`
        sleep 1

    done

}

start_server_old() {
    if [ -r $pid_file ] ; then
        echo "sys-monitor is running"
        exit 1
    fi

    echo -n "Starting sys-monitor node"

    $exec_prefix --node --pprof --cross --config $conf_file --pid $pid_file &

    if [ "$?" != 0 ] ; then
        echo " failed"
        exit 1
    fi

    wait_for_pid created $pid_file

    if [ -n "$try" ] ; then
        echo " failed"
        exit 1
    else
        echo " done"
    fi
}

start_server() {
    echo -n "Starting nohup node mode sys-monitor "

    nohup $exec_prefix --node --pprof --cross --config $conf_file >> ./out.log 2>&1 &

    if [ "$?" != 0 ] ; then
        echo " failed"
        exit 1
    fi

    wait_for_pid created $pid_file

    if [ -n "$try" ] ; then
        echo " failed"
        exit 1
    else
        echo " done"
    fi
}

stop_server() {
    echo -n "Gracefully shutting down sys-monitor node"

    if [ ! -r $pid_file ] ; then
        echo "warning, no pid file found - sys-monitor node is not running ?"
        exit 1
    fi

    kill -QUIT `cat $pid_file`

    wait_for_pid removed $pid_file

    if [ -n "$try" ] ; then
        echo " failed. Use force-quit"
        exit 1
    else
        echo " done"
    fi;
}

case "$1" in
    start)
        start_server
    ;;

    stop)
        stop_server
    ;;

    start_old)
        start_server_old
    ;;

    restart)
        stop_server
        start_server
    ;;
esac