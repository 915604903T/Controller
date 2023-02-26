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
    mesh1 = o3d.io.read_triangle_mesh(args.file1)
    mesh1.remove_duplicated_vertices()
    mesh1.remove_degenerate_triangles()
    mesh2 = o3d.io.read_triangle_mesh(args.file2)
    mesh2.remove_duplicated_vertices()
    mesh2.remove_degenerate_triangles()
    mesh2.transform(pose)
    mesh = mesh1 + mesh2
    o3d.io.write_triangle_mesh(args.outFile, mesh)
    print("success merge", args.file1, args.file2, "to", args.outFile)
    os.remove(args.file1)
    os.remove(args.file2)
    print("remove merged ply files")


