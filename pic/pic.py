from cProfile import label
import linecache
from matplotlib.font_manager import FontProperties
import matplotlib.pyplot as plt
from result import *
import numpy as np
import os
# # ThemiX
# x = [2567.101172,5051.544463,9488.566278,17540.78232,41948.75261,57868.36105]
# y = [10440.9, 10610.27, 10523.88, 11391.75, 14296.58,17279.15]
# # Honey Badger BFT
# x1 = [1.65808018, 8048.726002, 28186.3939, 77809.81711, 99480.56931]
# y1 = [4017.37, 4200.37, 4781.86, 8450.71, 15538.87]
# #HotStuff
# x2 = [9762, 19435, 48506, 95198, 95706]
# y3 = [15137, 14549, 15429, 15169,17153]

def pic(x, y, bottom, top, label, name):
    plt.plot(x, y, linewidth=2.5, label=label, color='limegreen', marker='s', markerfacecolor='none', markersize=10, linestyle='--')
    plt.xlabel("Throughput(tx/s)", fontsize=16)
    plt.ylabel("Average Latency(s)", fontsize=16)
    plt.ylim(bottom, top)
    plt.legend(fontsize="x-large")
    plt.tick_params(labelsize=13)
    plt.show()
    plt.close()

def pic2(x,y,x1,y1,name):
    plt.plot(x,y,linewidth=2.5,label='Dumbo 4 nodes',color='green',marker='o',markerfacecolor='none', markersize=10, linestyle=':')
    plt.plot(x1,y1,linewidth=2.5,label='Dumbo 7 nodes',color='purple',marker='^',markerfacecolor='none', markersize=10, linestyle='--')
    # plt.plot(x2,y3,linewidth=2,label='HotStuff',color='black',marker='v',markerfacecolor='none',linestyle='-.')
    # plt.scatter(x,y,marker='o',c='',edgecolors='black',label='list1')
    # plt.scatter(x1,y1,marker='^',c='',edgecolors='black',label='list2')
    # plt.title("HBBFT:105 Nodes in 5 areas,Optimal efficiency", fontsize=10)
    plt.xlabel("Throughput(ktx/sec)", fontsize=16)
    plt.ylabel("Average Latency(s)", fontsize=16)
    plt.ylim(0)
    plt.legend(fontsize="x-large")
    plt.tick_params(labelsize=13)
    plt.savefig(name)
    plt.close() 
    # plt.show()

def pic3(x,y,x1,y1,x3,y3,name):
    plt.plot(x,y,linewidth=2.5,label='ThemiX',color='b',marker='o',markerfacecolor='none', markersize=10, linestyle=':')
    plt.plot(x3,y3,linewidth=2.5,label='ThemiX-rs',color='g',marker='v',markerfacecolor='none', markersize=10, linestyle='-.')
    plt.plot(x1,y1,linewidth=2.5,label='ACS',color='r',marker='^',markerfacecolor='none', markersize=10, linestyle='--')

    plt.xlabel("Throughput(ktx/sec)", fontsize=16)
    plt.ylabel("Average Latency(s)", fontsize=16)
    plt.ylim(0)
    plt.legend()
    plt.tick_params(labelsize=13)
    plt.savefig(name)
    plt.close() 
    # plt.show()
# def pic(x,y,x1,y1,x3,y3,name):
#     plt.plot(x,y,linewidth=2,label='ThemiX',color='r',marker='o',markerfacecolor='none', markersize=10, linestyle=':')
#     plt.plot(x1,y1,linewidth=2,label='ACS',color='g',marker='^',markerfacecolor='none', markersize=10, linestyle='--')
#     plt.plot(x3,y3,linewidth=2,label='ThemiX-rs',color='b',marker='v',markerfacecolor='none', markersize=10, linestyle='-.')
#     plt.xlabel("Throughput(ktx/sec)", fontsize=16)
#     plt.ylabel("Average Latency(s)", fontsize=16)
#     plt.ylim(0)
#     plt.legend()
#     plt.tick_params(labelsize=13)
#     plt.savefig(name)
#     plt.close() 
    # plt.show()

def pic4(x,y,x1,y1,x3,y3,x4,y4,name):
    plt.plot(x,y,linewidth=2.5,label='ThemiX',color='b',marker='o',markerfacecolor='none', markersize=10, linestyle=':')
    plt.plot(x3,y3,linewidth=2.5,label='ThemiX-rs',color='g',marker='v',markerfacecolor='none', markersize=10, linestyle='-.')
    plt.plot(x1,y1,linewidth=2.5,label='ACS',color='r',marker='^',markerfacecolor='none', markersize=10, linestyle='--')
    plt.plot(x4, y4, linewidth=2.5, label='Dumbo', color='purple', marker='*', markerfacecolor='none', markersize=10, linestyle='-')

    plt.xlabel("Throughput(ktx/sec)", fontsize=16)
    plt.ylabel("Average Latency(s)", fontsize=16)
    plt.ylim(0)
    plt.legend()
    plt.tick_params(labelsize=13)
    plt.savefig(name)
    plt.close() 

if not os.path.exists("./graphs"):
    os.mkdir("./graphs")

pic(normal_3_x, normal_3_y, 8, 16, 'ThemiX (3 nodes)', './graphs/normal_3.jpg')
pic(normal_7_x, normal_7_y, 16, 25, 'ThemiX (7 nodes)', './graphs/normal_7.jpg')
# pic2(d4_x, d4_y, d7_x, d7_y, "./graphs/dumbo.jpg")
pic3(tm_1_x1,tm_1_y1,hb_x1,hb_y1,tm_2_x1,tm_2_y1,"./graphs/3nodes_0crashes.pdf")
pic4(tm_1_x1,tm_1_y1,hb_x1,hb_y1,tm_2_x1,tm_2_y1,d4_x, d4_y, "./graphs/4nodes_0crashes.pdf")
pic3(tm_1_x3,tm_1_y3,hb_x3,hb_y3,tm_2_x3,tm_2_y3,"./graphs/67nodes_0crashes.pdf")

# pic3(tmc_x1,tmc_y1,hbc_x1,hbc_y1,tmcs_x1,tmcs_y1,"./graphs/3nodes_1crashes.pdf")
# pic3(tmc_x2,tmc_y2,hbc_x2,hbc_y2,tmcs_x2,tmcs_y2,"./graphs/7nodes_1crashes.pdf")
# pic3(tmc_x3,tmc_y3,hbc_x3,hbc_y3,tmcs_x3,tmcs_y3,"./graphs/7nodes_3crashes.pdf")


plt.plot(tm_1_x2,tm_1_y2,linewidth=2.5,label='ThemiX(7 nodes)',color='b',marker='o',markerfacecolor='none', markersize=10, linestyle='--')
plt.plot(tm_2_x2,tm_2_y2,linewidth=2.5,label='ThemiX-rs(7 nodes)',color='g',marker='v',markerfacecolor='none', markersize=10,linestyle='--')
plt.plot(hb_1_x2,hb_1_y2,linewidth=2.5,label='ACS(7 nodes)',color='r',marker='^',markerfacecolor='none', markersize=10,linestyle='-.')
plt.plot(hb_2_x2,hb_2_y2,linewidth=2.5,label='ACS(10 nodes)',color='Grey',marker='x',markerfacecolor='none',markersize=10,linestyle='-.')
# plt.plot(d7_x, d7_y, linewidth=2.5,label='Dumbo(7 nodes)',color='Purple',marker='*',markerfacecolor='none',markersize=10,linestyle='-')
plt.xlabel("Throughput(ktx/sec)", fontsize=16)
plt.ylabel("Average Latency(s)", fontsize=16)
plt.ylim(0)
plt.legend()
plt.tick_params(labelsize=13)
plt.savefig("./graphs/7nodes.pdf")
plt.close()


width = 0.3
x = np.arange(len(name_list))

fig = plt.figure(figsize=(12, 9))
ax1 = fig.add_subplot(111)
ax1.set_xlabel('Nodes', fontsize=12)
ax1.set_ylabel('Throughput (ktx/s)', fontsize=16)
ax1.tick_params(labelsize=13)
ax1.bar(x - width/2, throughput_list_without_sign, tick_label=name_list, width=width,label=u"Throughput(ThemiX)", fill = False, hatch = "/")
ax1.bar(x + width/2, throughput_list_with_sign, tick_label=name_list, width=width,label=u"Throughput(ThemiX-rs)",fill = False, hatch = "X")
ax1.set_xticks(x)
ax1.set_xticklabels(name_list)
ax2 = ax1.twinx()
ax2.set_ylabel('Latency (s)', fontsize=16)
ax2.tick_params(labelsize=13)
ax2.plot(x, latency_list_without_sign, linewidth=3, label=u'Latency(ThemiX)',color='b',marker='o',markerfacecolor='none', markersize=10, linestyle='--')
ax2.plot(x, latency_list_with_sign,  linewidth=3, label=u'Latency(ThemiX-rs)',color='g',marker='v',markerfacecolor='none', markersize=10, linestyle='-.')

lines = []
labels = []
for ax in fig.axes:
    axLine, axLabel = ax.get_legend_handles_labels()
    lines.extend(axLine)
    labels.extend(axLabel)

fig.legend(lines, labels,fontsize="xx-large", loc = 'center',bbox_to_anchor=(0.30,0.78))
plt.savefig("./graphs/nodes.pdf")
plt.close()
