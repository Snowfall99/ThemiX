import sys

nodes = int(sys.argv[1])
name = str(sys.argv[2])

for i in range(nodes):
  filename = 'log/client' + str(i) + '/' + name
  latencies = []
  file = open(filename, 'r')
  lines = file.readlines()
  file.close()
  for line in lines:
    try:
      latency = float(line)
      latencies.append(latency)
    except:
      break
  latencies.sort()
  # print(len(latencies))
  print(latencies[int(0.95*len(latencies))])
