import boto3
import json
import sys

regions = ["us-east-1","us-west-1","ap-northeast-1","ap-southeast-2","ap-southeast-1","ap-east-1","eu-central-1","eu-west-1","eu-west-3","me-south-1","sa-east-1","ca-central-1"]
access_key = ""
secret_key = ""

total = {}
total["nodes"] = []
clients = {}
clients["nodes"] = []
server_id = 0
client_id = 0
for region in regions:
    print("region:",region)
    ec2 = boto3.client('ec2',aws_access_key_id=access_key, aws_secret_access_key=secret_key,region_name=region)
    Tags = [{'Key': 'Name', 'Value': 'Free'}]
    Filter = [
        {
            'Name': 'key-name',
            'Values': [
                'aws',
            ]
        }
    ]
    response = ec2.describe_instances(Filters=Filter)
    instances = []
    for i in range(len(response['Reservations'])):
        instances += response['Reservations'][i]['Instances']
    print(len(instances))

    # --------------------------------
    # --------nodes-----------------
    for i in range(len(instances)):
        status = instances[i]['State']['Name']
        if status != "running":
            continue
        instance = {}
        instance['Id'] = server_id
        server_id += 1
        instance['InstanceId'] = instances[i]['InstanceId']
        instance['InstanceType'] = instances[i]['InstanceType']
        instance['PublicIpAddress'] = instances[i]['PublicIpAddress']
        instance['PrivateIpAddress'] = instances[i]['PrivateIpAddress']
        instance['ServerURL'] = "http://" + instances[i]['PublicIpAddress'] +":6000/client"

        total['nodes'].append(instance)





print("----- begin to load----")
file = "./nodes.json"
with open(file,"w") as f:
    json.dump(total,f)
print("----- load success ----")


for item in range(len(total['nodes'])):
    total['nodes'][item]['ServerURL'] = "http://" + total['nodes'][item]['PublicIpAddress'] +":6100/client"

print("----- begin to load ----")
file = "../client/clients.json"
with open(file,"w") as f:
    json.dump(total,f)
print("----- load success ----")
