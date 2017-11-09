## The intuitionistic fuzzy C-means algorithm for Golang

An implementation of IFCM algorithm based on the paper `A novel intuitionistic fuzzy C means clustering algorithm and its application to medical images`. 

## The IFCM struct
```golang
type IFCM struct {
    data []float64       // the input sequence values
    c    int             // clustering size
    m    int             // parameter m, default 2
    t    float64         // the termination condition, default 0.005
    a    float64         // parameter a, default 0.85
    max  int             // maximum number of iterations, default 100
    v    []float64       // keeping the clustering center V
    u    [][]float64     // keeping the membership degree U
}
```

## Usage

```golang
import{
    "mervin.me/IFCM/ifcm"
}
var data []float64
var c int = 3
ifcm := ifcm.NewIFCM(data, c)
ifcm.Run()
```

