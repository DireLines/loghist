import sys
import json
from matplotlib import pyplot as plt
import time

filters = sys.argv[1:]
def filterDict(predicate,dictObj):
    result = {}
    for (key, value) in dictObj.items():
        if predicate(key, value):
            result[key] = value
    return result
def filters_contain_key(key,value):
    for filter in filters:
        if filter in key:
            return True
    return False
data = {}
plt.ion()
plt.show()
while True:
    batch = json.loads(sys.stdin.readline())
    if filters:
        batch = filterDict(filters_contain_key,batch)
    for k in batch:
        if k not in data:
            data[k] = []
        data[k].extend(batch[k])
    if not data:
        continue
    nums = data.values()
    plt.cla()
    plt.hist(nums,30,stacked=True)
    plt.legend(data.keys())
    plt.draw()
    plt.pause(0.001)

    