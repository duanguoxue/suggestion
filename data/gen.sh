awk -F " " '{print $1 "###" srand(rand()*1000000000) }' SogouLabDic.dic  > dict.txt
