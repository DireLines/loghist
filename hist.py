import sys
import asyncio
import subprocess
from matplotlib import pyplot as plt

indicator = sys.argv[1]
lines = []
while True:
    num_lines = 5000
    lines.extend(sys.stdin.readlines(num_lines))
    plt.hist(
        list(
            map(float,#convert to number
                map(lambda s: s.split()[-2],#extract data point
                    filter(lambda s: indicator in s, lines)
                )
            )
        )
    )
    plt.show()