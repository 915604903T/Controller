import open3d as o3d
import numpy as np
import os
import argparse

def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f1", "--file1", dest='file1', required=True)
    parser.add_argument("-f2", "--file2", dest='file2', required=True)
    parser.add_argument("-p", "--pose", dest='poseFile', required=True)
    parser.add_argument("-o", "--output", dest='outFile', required=True)
    args = parser.parse_args()
    return args

if __name__=='__main__':
    args = parse_args()
    pose = np.empty([4, 4], dtype='float64')
    poseFile = open(args.poseFile)
    for i in range(4):
        line = poseFile.readline()
        pose[i] = np.fromstring(line, dtype='float64', sep=' ')
    pc1 = o3d.io.read_point_cloud(args.file1)
    pc2 = o3d.io.read_point_cloud(args.file2)
    pc2.transform(pose)
    pc = pc1 + pc2
    pc = pc.voxel_down_sample(voxel_size=0.00001)
    pc.estimate_normals()
    o3d.io.write_point_cloud(args.outFile, pc)
    print("success merge", args.file1, args.file2, "to", args.outFile)


