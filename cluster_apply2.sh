for d in $(ls yaml_parser/clusters/) ; do kubectl apply -f "yaml_parser/clusters/${d}/"; done
