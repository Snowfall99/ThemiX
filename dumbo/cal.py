import argparse
import re

if __name__ == '__main__':
  parser = argparse.ArgumentParser()
  parser.add_argument('--name', metavar='name', required=True,
                      help='name of log file', type=str)
  parser.add_argument('--N', metavar='N', required=True,
                      help='number of nodes', type=int)
  parser.add_argument('--K', metavar='K', required=True,
                      help='rounds to execute', type=int)
  args = parser.parse_args()

  name = args.name
  N = args.N
  K = args.K

  count = 0
  latency = 0
  throughput = 0
  for i in range(N):
    filename = './log/' + name + '-' + str(i) + '.log'
    file = open(filename, 'r+')
    lines = file.readlines()
    file.close()
    result = lines[len(lines)-1]
    try:
      searchObj = re.search(r'(.*) breaks in (.*) seconds with total delivered Txs (.*)', result, re.M|re.I)
      time = searchObj.group(2)
      txs = searchObj.group(3)
      latency += float(time)/K
      throughput += float(txs)/float(time)
      count += 1
    except:
      continue

  print("latency(s): ", latency/count)
  print("throughput(tps): ", throughput)
