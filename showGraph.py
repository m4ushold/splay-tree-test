import matplotlib.pyplot as plt

name_map = dict()

def plot_results(name, results):
    bench_name = []
    max_heights = []
    max_rotations = []
    sum_rotations = []
    times = []

    for result in results:
        bench_name.append(result['bench-name'])
        max_heights.append(result['max-height'])
        max_rotations.append(result['max-rotation'])
        sum_rotations.append(result['sum-rotation'])
        times.append(result['ns/op'])

    fig, ax1 = plt.subplots(figsize=(12, 8))    

    bench_name = [name_map[i] for i in bench_name]

    # max height
    ax1.plot(bench_name, max_heights, marker='o', color='#1f77b4', label='max height')
    ax1.set_xlabel('Splay trees')
    ax1.set_ylabel('max heights', color='#1f77b4')
    ax1.tick_params(axis='y', labelcolor='#1f77b4')

    # sum rotations
    ax2 = ax1.twinx()
    ax2.plot(bench_name, sum_rotations, marker='^', color='#2ca02c', label='count of rotation calls')
    ax2.set_ylabel('rotation counts', color='#2ca02c')
    ax2.tick_params(axis='y', labelcolor='#2ca02c')

    # time
    ax3 = ax1.twinx()
    ax3.spines['right'].set_position(('outward', 60))
    ax3.plot(bench_name, times, marker='x', color='#9467bd', label='ns/op')
    ax3.set_ylabel('ns / op', color='#9467bd')
    ax3.tick_params(axis='y', labelcolor='#9467bd')

    # max rotaion
    ax4 = ax1.twinx()
    ax4.spines['right'].set_position(('outward', 120))
    ax4.plot(bench_name, max_rotations, marker='x', color='#ff7f0e', label='ns/op')
    ax4.set_ylabel('rotation', color='#ff7f0e')
    ax4.tick_params(axis='y', labelcolor='#ff7f0e')

    plt.title(name[:-4])
    
    fig.tight_layout()
    ax1.legend(loc='upper left')
    ax2.legend(loc='upper center')
    ax3.legend(loc='upper right')
    ax4.legend(loc='upper right')

    plt.savefig('result/'+name[:-4])


def parse_bench_file(filename):
    res=[{"benchmark":"", "results":[], "ns/op":""}]
    with open(filename, 'r') as f:
        for line in f:
            if line[:9] == "Benchmark":
                res.append({
                    "benchmark":line,
                    "results":[],
                    "ns/op": ""
                })
            elif line[:6] == 'result':
                res[-1]['results'].append(line)
            elif line[:8] == '10000000':
                res[-1]['ns/op'] = line
    bench_result = []
    for result in res:
        if result['benchmark'] != '' and len(result['results'])>0:
            mxH,mxR,sumR=0,0,0

            if len(result['results']) == 8:
                result['results'] = result['results'][:-1]
                
            for re in result['results']:
                h,r,sr=map(int,re[7:-2].split(','))
                mxH+=h
                mxR+=r
                sumR+=sr
            mxH/=len(result['results'])
            mxR/=len(result['results'])
            sumR/=len(result['results'])
                
            bench_result.append({
                'bench-name': result['benchmark'].strip(),
                'max-height': mxH, 
                'max-rotation': mxR,
                'sum-rotation': sumR,
                'ns/op': float(result['ns/op'].split()[1])
            })
    return bench_result

if __name__ == "__main__":
    defulat_res = parse_bench_file('result/basic-tree-result.txt')
    with open('result/splay-tree-numbers.txt', 'r') as f:
        for line in f:
            n, nm=line.split()
            name_map[nm]=n

    for file in ['random-result.txt', 'max-height-splay-result.txt', 'stlb-result.txt', 'result.txt']:
        res = parse_bench_file('result/'+file)
        res = defulat_res + res
        plot_results(file, res)

                