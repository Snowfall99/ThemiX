import sys

nodes = int(sys.argv[1])
name = str(sys.argv[2])

for i in range(nodes):
  filename = 'log/client' + str(i) + '/' + name
  txs = 0
  latency = 0
  file = open(filename, 'r')
  lines = file.readlines()
  file.close()
  for line in lines:
    try:
      latency += float(line)
      txs += 1
    except:
      break
  print('Latency(ms): ', latency/txs)
