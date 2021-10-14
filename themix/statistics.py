from numpy import mean

f = open("coordinator.txt")
line = f.readline()

list = []
while line:
    if line.find("start") != -1:
        time = line.split(":")[1].strip()
        start = time
    else:
        time = line.split(":")[1].strip()
        list.append(int(time) - int(start))
    line = f.readline()

list.sort()

print("mean: ", mean(list))
print("fastest: ", list[0])
print("slowest: ", list[len(list)-1])