def div1000(num):
    res = []
    for i in range(len(num)):
        res.append(num[i]/1000)
    return res

def div3000(num):
    res = []
    for i in range(len(num)):
        res.append(num[i]/3000)
    return res
# 3nodes:f = 0


# ThemiX Without Signature
tm_1_x1 = div1000([9.70384373, 16713.9233, 24561.86936, 31385.47762, 35928.14371, 36304.06507])
tm_1_y1 = div1000([308.1450216, 364.5939394, 499.8060942, 793.0393013, 1045.373563, 1316.442029])

# Honey Badger BFT
hb_x1 = div1000([4.880946488, 9793.620269, 16144.10402, 25150.24694, 35191.19431, 41269.12382, 43740.33701, 45991.55297, 45903.14437])
hb_y1 = div1000([835.6424116, 832.4419087, 1019.454315, 1320.596721, 1906.521327, 2480.527607, 2955.583942, 3213.826772, 3785.805556])

# ThemiX With Signature
tm_2_x1 = div1000([9.962763757, 13468.68365, 19233.62911, 22268.2462, 20976.62185])
tm_2_y1 = div1000([304.4249578, 454.488665, 642.0035461, 1136.339623, 1861.459184])


# 7nodes:f = 0
# ThemiX Without Signature
tm_1_x2 = div3000([69.8203, 59201.184, 115316.1751, 155633.4941,154057.0976])
tm_1_y2 = div1000([495.88, 583.52, 1192.31, 2232.67,3305.34])

# ThemiX With Signature
tm_2_x2 = div1000([14.09729444, 20033.27897, 24739.44134, 30637.91117, 35067.43738, 34819.93799])
tm_2_y2 = div1000([507.0972389, 720.3441227, 1172.055851, 1988.934884, 2656.074534, 3454.264])

# Honey Badger BFT 7 nodes
hb_1_x2 = div1000([5.717272698, 10790.89783, 17829.06786, 28736.5497, 39629.3565, 45790.48843, 48048.19355, 52000.76099, 53207.54717])
hb_1_y2 = div1000([1289.62614, 1372.932258, 1680.165354, 2104.651961, 2315.227027, 2531.568047, 2778.350649, 3022.41958, 3577.166667])

# Honey Badger BFT 10 nodes
hb_2_x2 = div1000([6.652245336, 11663.31886, 21434.27203, 36169.07089, 48707.99175, 56117.95775, 65670.97566, 70277.36132, 70013.46413])
hb_2_y2 = div1000([1599.960422, 1850.169697, 2032.926421, 2455.673387, 2732.767857, 3011.649038, 3508.556818, 4125.28, 4594.875969])


# 100nodes:f=0
# ThemiX Without Signature
tm_1_x3 = div3000([10276.5681, 325482.7881, 392626.6054, 414118.4044, 420344.9087, 421095.8262]) 
tm_1_y3 = div1000([1190.42, 1618.49, 2620.39, 2999.45, 4941.04, 6205.79])


# ThemiX Without Signature
tm_2_x3 = div1000([49.89594688, 3473.631083, 20207.5862, 29279.41264, 29457.40974]) 
tm_2_y3 = div1000([1159.330298, 1574.133392, 2872.108173, 3170.023605, 4788.599398])

# Honey Badger BFT
hb_x3 = div3000([11396.94, 110068.4799, 166005.1308, 238958.7091, 281603.5883, 306485.2274])
hb_y3 = div1000([19454.58, 20631.72, 26180.98, 28028.31, 31385.07, 34548.25])



# ------------------------------ Crash ------------------------------
# 3nodes:f = 1
# ThemiX Without Signature
tmc_x1 = div3000([1.424921885, 1297.479646, 5070.779632, 21664.24749, 28655.42383])
tmc_y1 = div1000([7014,7704.96,7886.17,9229.05,10468.44])

# ThemiX With Signature
tmcs_x1 = div1000([0.283730872, 1907.044144, 3458.789168, 4516.141887, 6311.52968, 8216.504904])
tmcs_y1 = div1000([8806.916667, 9960.875, 10401.8, 10745.75, 11535.82143, 11977.69231])

# Honey Badger BFT
hbc_x1 = div1000([2.504458326, 24169.18429, 40066.8102, 46534.17353, 62808.69043880435, 65760.03068801432])
hbc_y1 = div1000([1225.508065, 2093.751724, 4152.797297, 6020.352941, 13035.875, 15535.761904761905])

# 7nodes:f=1
# ThemiX Without Signature
# tmc_x2 = div3000([35.1031154, 32590.60187, 46331.81551, 76124.98012, 91451.29225])
# tmc_y2 = div1000([845.05, 909.18, 1271.91, 1542.26, 3141.41])
tmc_x2 = div1000([7, 13660.3, 47085.75, 54010, 49195.5])
tmc_y2 = div1000([767, 785.145,	1790.63, 2599.45, 3331.55])


# ThemiX With Signature
# tmcs_x2 = div1000([7.026955428, 11362.5077, 16687.80541, 22429.26836, 25393.15357, 24603.7524])
# tmcs_y2 = div1000([865.5867689, 1075.594048, 1460.020635, 2220.264059, 2985.026316, 3606.796791])
tmcs_x2 = div1000([7.026955428, 11362.5077, 16687.80541, 25393.15357, 24603.7524])
tmcs_y2 = div1000([865.5867689, 1075.594048, 1460.020635, 2985.026316, 3606.796791])

# Honey Badger BFT
hbc_x2 = div3000([4.55503573, 20695.07871, 69293.20927, 136430.9282, 155141.844])
hbc_y2 = div1000([1922.69, 2081.79, 2508.37,6134.07,11080.61])

# 7nodes:f=3
# ThemiX Without Signature
tmc_x3 = div3000([2.626740215, 2325.704318, 19992.4473, 45813.18402, 61680.4214])
tmc_y3 = div1000([7611.96, 8595.3, 9997.72, 13095.71, 16205.7])

# ThemiX With Signature
tmcs_x3 = div1000([0.496488649, 852.5814606, 3397.748991, 6006.006006, 9278.560957])
tmcs_y3 = div1000([9519.0625, 10966.72881, 11372.98214, 12230.15385, 13405.20833])

# Honey Badger BFT)
hbc_x3 = div3000([1.65808018, 8048.726002, 28186.3939, 77809.81711, 99480.56931])
hbc_y3 = div1000([4017.37, 4200.37, 4781.86, 8450.71, 15538.87])

name_list = ['3', '5', '11', '31', '67']
# latency
latency_list_without_sign = div1000([308, 403, 530, 1013, 1146.61])
latency_list_with_sign = div1000([454, 508, 903, 1333, 1574])
# throughput
# throughput_list_without_sign = div3000([58911, 61434, 106974, 251039, 333297])
throughput_list_without_sign = div3000([58911, 37732.74*3, 56743.08*3, 251039, 333297])
throughput_list_with_sign = div1000([19233, 30800, 35067, 33884, 29279])



# timeout data
tm_delta_x1 = div1000([0.849919448, 1435.397662, 2849.773403, 5170.965031, 8448.108136, 12938.03329, 14646.55549, 21479.410683077804, 22420.495742072522])
tm_delta_y1 = div1000([5044.319672, 6063.490196, 6282.104167, 6992.738636, 8165.026316, 9044.323529, 9680.390625, 16133.697368421053, 19001.515625])

tm_delta_x2 = div1000([0.496488649, 852.5814606, 3397.748991, 6006.006006, 9278.560957, 18276.025701863513, 16272.224748139462])
tm_delta_y2 = div1000([9519.0625, 10966.72881, 11372.98214, 12230.15385, 13405.20833, 19436.6875, 20505.75])

tm_delta_x3 = div1000([0.341067351, 640.5365436, 1259.84872, 2441.350372, 4350.242889, 6949.925785, 9283.804141,11520.700458587882,12291.956984723793,14061.073413436914,15551.24296549244,13557.339571481736])
tm_delta_y3 = div1000([14026.35, 15110.35714, 15581.55769, 16086.55769, 17167.77083, 18248.40909, 18990.02941,20774.866666666665,21664.053571428572,22390.89285714286,23683.46153846154, 24838.583333333332])

d4_x = div1000([6.674, 7707.24, 14241.719, 24363.035, 28404.76955])
d4_y = [1.79, 1.94, 2.10, 2.46, 3.17]

d7_x = div1000([14.759, 10096, 18798.527, 33839.056, 43814.52, 54065, 58741, 63624])
d7_y = [2.372, 2.475, 2.66, 2.96, 3.42, 3.69, 4.26, 4.72]

normal_3_x = [0.425, 1750, 3750, 4666.7, 6750, 8666.7, 19500, 21000]
normal_3_y = [8.736, 10.01, 10.329, 11.086, 11.948, 12.362, 12.699, 14.204]

normal_7_x = [0.64, 3208, 6417, 12833, 18666, 28000, 34222]
normal_7_y = [17.01, 18.85, 18.16, 18.77, 20.06, 20.94, 23.49]
