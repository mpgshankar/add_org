export EXISTMSP=()
# echo $EXISTMSP
for ((i=1; i <= 2; i++))
{
    EXISTMSP+=('test'$i'MSP')
}
# printf "\'%s\'," "${EXISTMSP[@]}"

# export MSP=""
# for index in "${!EXISTMSP[@]}"
# do
# if [$index == ${!EXISTMSP[@]} ] then
#         MSP+=${EXISTMSP[index]}
# else
#         MSP+=${EXISTMSP[index]},
# fi    
#     # echo "$index ${EXISTMsSP[index]}"
# done
# echo $MSP



export CLINODEID=$(docker node ls -f name=len --format "{{.ID}}")

echo $CLINODEID