for i in yaml_parser/clusters/*; do kubectl apply -f $(echo $i); done
