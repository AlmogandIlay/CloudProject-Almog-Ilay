



path = "C:\\Users\\אילאי\\OneDrive\\Desktop\\sample-2mb-text-file.txt"




f = open(path, "r")


a = f.readlines()

for i in a:
    print(i)