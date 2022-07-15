import matplotlib.pyplot as plt

normal_3_x = [0.425, 1750, 3750, 6500, 9000, 12000, 19500, 21000]
normal_3_y = [8.736, 10.01, 10.329, 11.086, 11.948, 12.362, 12.699, 14.204]

sign_3_x = [0.37, 1750, 3250, 6000, 9000, 12000, 15000, 17111]
sign_3_y = [11.14, 11.68, 12.26, 13.49, 14.26, 14.62, 15.57, 16.57]

normal_7_x = [0.64, 3208, 6417, 12833, 18666, 28000, 34222]
normal_7_y = [17.01, 18.85, 18.16, 18.77, 20.06, 20.94, 23.49]

sign_7_x = [0.58, 2916, 5833, 10500, 18666, 24111]
sign_7_y = [17.26, 18.26, 18.78, 20.997, 21.188, 25.79]

delta_1_x = [1.28, 4750, 8000, 14000, 19833.33, 30333.33, 35000]
delta_1_y = [6.45, 7.57, 8.97, 9.33, 8.77, 12.2, 14.72]

delta_2_x = [0.64, 3208, 6417, 12833, 18666, 28000, 34222]
delta_2_y = [17.01, 18.85, 18.16, 18.77, 20.06, 20.94, 23.49]

delta_3_x = [0.39, 3888.89, 11000, 14666.67, 21000, 30333.33, 36400]
delta_3_y = [26.68, 27.88, 28.02, 28.56, 28.23, 31.21, 33.55]

regions = ['N.Virginia', 'N.California', 'Tokyo', 'Sydney', 'Hong Kong', 'Frankfurt', 'Ireland']
latency = [17144.125, 17079.25, 16941.0, 16918.25, 16915.375, 17008.75, 17067.0]

crash_3_x = [0.26, 2000, 4000, 7333.33, 10000, 16666.67]
crash_3_y = [12.09, 13.84, 14.26, 15.16, 15.09, 16.17]

crash_rs_3_x = [0.22, 2000, 3666.67, 6666.67, 10000, 14444.44]
crash_rs_3_y = [12.61, 14.11, 15.14, 15.09, 16.29, 18.87]

crash_7_x = [0.26, 1333.33, 4444.44, 6666.67, 8888.89]
crash_7_y = [26.10, 27.93, 28.18, 29.11, 30.46]

crash_rs_7_x = [0.26, 1333.33, 4444.44, 6000, 7333.33]
crash_rs_7_y = [26.11, 27.41, 28.84, 29.94, 30.31]

plt.plot(normal_3_x, normal_3_y, label='ThemiX (3 nodes)', color='darkorange', linewidth=3, marker='s', markersize=10, markerfacecolor='none', linestyle='-.')
plt.plot(sign_3_x, sign_3_y, label='ThemiX-rs (3 nodes)', color='olivedrab', linewidth=3, marker='^', markersize=10, markerfacecolor='none', linestyle='-.')
plt.xlabel("Throughput(tx/s)", fontsize=12, fontweight='bold')
plt.ylabel("Average Latency(s)", fontsize=12, fontweight='bold')
plt.ylim(8, 18)
# plt.grid('both', linestyle='--')
plt.legend(fontsize="large")
plt.tick_params(labelsize=13)
plt.savefig("./graphs/normal_3_fault_free.pdf", bbox_inches='tight')
plt.close()

plt.plot(normal_7_x, normal_7_y, label='ThemiX (7 nodes)', color='darkorange', linewidth=3, marker='s', markersize=10, markerfacecolor='none', linestyle='-.')
plt.plot(sign_7_x, sign_7_y, label='ThemiX-rs (7 nodes)', color='olivedrab', linewidth=3, marker='^', markersize=10, markerfacecolor='none', linestyle='-.')
plt.xlabel("Throughput(tx/s)", fontsize=12, fontweight='bold')
plt.ylabel("Average Latency(s)", fontsize=12, fontweight='bold')
plt.ylim(16, 27)
# plt.grid('both', linestyle='--')
plt.legend(fontsize="large")
plt.tick_params(labelsize=13)
plt.savefig("./graphs/normal_7_fault_free.pdf", bbox_inches='tight')
plt.close()

plt.plot(delta_1_x, delta_1_y, label='$\Delta$=500ms', color='slategray', linewidth=2, marker='x', markersize=12, markerfacecolor='none', linestyle='-.')
plt.plot(delta_2_x, delta_2_y, label='$\Delta$=1250ms', color='orangered', linewidth=2, marker='*', markersize=12, markerfacecolor='none', linestyle='-.')
plt.plot(delta_3_x, delta_3_y, label='$\Delta$=2000ms', color='darkcyan', linewidth=2, marker='+', markersize=12, markerfacecolor='none', linestyle='-.')
plt.xlabel("Throughput(tx/s)", fontsize=12, fontweight='bold')
plt.ylabel("Average Latency(s)", fontsize=12, fontweight='bold')
plt.ylim(0)
# plt.grid('both', linestyle='--')
plt.legend(fontsize="large")
plt.tick_params(labelsize=13)
plt.savefig("./graphs/delta.pdf", bbox_inches='tight')
plt.close()

plt.figure(figsize=(10, 6))
plt.bar(regions, latency, label='ThemiX (7 nodes)', width=0.6, color='darkturquoise', hatch='/', edgecolor='black')
plt.xlabel('AWS EC2 Regions', fontsize=12, fontweight='bold')
plt.ylabel('Average Latency(ms)', fontsize=12, fontweight='bold')
# plt.yticks(np.arange(20000, 21500, step=500))
plt.ylim(16500, 17200)
# plt.grid(axis='y', linestyle='-.', zorder=3, linewidth=0.5)
plt.legend()
plt.tick_params(labelsize=13)
plt.savefig('./graphs/regions.pdf', bbox_inches='tight')
plt.close()

plt.plot(crash_3_x, crash_3_y, label='ThemiX crash (3 nodes)', color='indigo', linewidth=3, marker='o', markersize=10, markerfacecolor='none', linestyle='-.')
plt.plot(crash_rs_3_x, crash_rs_3_y, label='ThemiX-rs crash (3 nodes)', color='firebrick', linewidth=3, marker='v', markersize=10, markerfacecolor='none', linestyle='-.')
plt.xlabel("Throughput(tx/s)", fontsize=12, fontweight='bold')
plt.ylabel("Average Latency(s)", fontsize=12, fontweight='bold')
plt.ylim(11, 20)
# plt.grid('both', linestyle='--')
plt.legend(fontsize="large")
plt.tick_params(labelsize=13)
plt.savefig("./graphs/normal_3_crash.pdf", bbox_inches='tight')
plt.close()

plt.plot(crash_7_x, crash_7_y, label='ThemiX crash (7 nodes)', color='indigo', linewidth=3, marker='o', markersize=10, markerfacecolor='none', linestyle='-.')
plt.plot(crash_rs_7_x, crash_rs_7_y, label='ThemiX-rs crash (7 nodes)', color='firebrick', linewidth=3, marker='v', markersize=10, markerfacecolor='none', linestyle='-.')
plt.xlabel("Throughput(tx/s)", fontsize=12, fontweight='bold')
plt.ylabel("Average Latency(s)", fontsize=12, fontweight='bold')
# plt.ylim(11, 20)
# plt.grid('both', linestyle='--')
plt.legend(fontsize="large")
plt.tick_params(labelsize=13)
plt.savefig("./graphs/normal_7_crash.pdf", bbox_inches='tight')
plt.close()
