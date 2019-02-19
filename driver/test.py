
import subprocess
import shlex
import os 
import signal 
import time 
import math 
import threading 
import psutil 

SAMPLES = "10000"
CPP_SERVER = "./cpp-server"
CPP_CLIENT = "./cpp-client -c {}".format(SAMPLES)
GO_SERVER = "./go-server"
GO_CLIENT = "./go-client -c {}".format(SAMPLES)

SERVER = CPP_SERVER
CLIENT = CPP_CLIENT


def percentile(sv, p):
    "Assuming sorted values." 
    assert 0 <= p <= 100
    index = min(len(sv) - 1, math.ceil(len(sv)*p/100))
    v = sv[index]
    return v


def print_results(server, client, svalues):
    print(server, 'server', '|', client, 'client', '|', 
            'p50:', percentile(svalues, 50), 'us', 
            'p90:', percentile(svalues, 90), 'us',
            'p95:', percentile(svalues, 95), 'us')


def monitor_thread(tag, pid, event):
    p = psutil.Process(pid)
    cpu = []
    mem = []
    while not event.is_set():
        try:
            cpu.append(p.cpu_percent(interval=1))
            mem.append(p.memory_percent())
        except psutil.NoSuchProcess as e:
            break
    scpu = sorted(cpu)
    smem = sorted(mem)
    print(tag, 'CPU: p50', percentile(scpu, 50), 'p90:', percentile(scpu, 90))
    print(tag, 'MEM: p50', percentile(smem, 50), 'p90:', percentile(smem, 90))

def launch_monitor_thread(tag, pid, event):
    threading.Thread(target=monitor_thread, args=(tag, pid, event)).start()

for lang, SERVER in {'cpp':CPP_SERVER, 'go':GO_SERVER}.items():
    with subprocess.Popen(shlex.split(SERVER), preexec_fn=os.setsid) as sp:
        # Wait for server to start ..
        time.sleep(5)
        for lang_client, CLIENT in {'cpp': CPP_CLIENT, 'go': GO_CLIENT}.items():
            with subprocess.Popen(shlex.split(CLIENT), stdout=subprocess.PIPE) as cp:
                event = threading.Event()
                launch_monitor_thread('server', sp.pid, event)
                output = cp.communicate(timeout=10)[0].decode("utf-8")
                event.set() 
                values = [int(line.strip()) for line in output.split()]
                svalues = sorted(values)
                print_results(lang, lang_client, svalues)
                time.sleep(5)
        print(sp.pid)
        os.killpg(os.getpgid(sp.pid), signal.SIGTERM)
