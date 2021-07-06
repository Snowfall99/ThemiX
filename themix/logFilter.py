for i in range(4):
    file0 = open("log/server" + str(i), "rb")
    lines = file0.readlines()
    file_filter0 = open("log_filter/server" + str(i), "wb")
    for line in lines:
        if '"round":0' in line: 
            file_filter0.write(line)

    file0.close()
    file_filter0.close()