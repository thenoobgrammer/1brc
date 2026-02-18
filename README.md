# 1brc

The one billion row is a challenge that surfaced on Jan 1st 2024. It's a popular programming competition amongst the dev community and it was designed to test the limits of modern Java.

### The objective

You are given an unsorted list of weather stations and the recorded temperature formatted as `<string: station name>;<double: measurement>` and the 
goal is simple: process one billion of these records in the least amount of time possible.
```csv
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
St. John's;15.2
Cracow;12.6
Bridgetown;26.9
Istanbul;6.2
Roseau;34.4
Conakry;31.2
Istanbul;23.0
```
The task is to write a Go (initially Java) program which reads the file, calculates the min, mean, and max temperature value per weather station, 
and emits the results on stdout like this
(i.e. sorted alphabetically by station name, and the result values per station in the format <min>/<mean>/<max>, rounded to one fractional digit):
```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
```


