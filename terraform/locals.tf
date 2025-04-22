locals {
  config_file_contents = file(var.config_file_local_path)

  controller_files = fileset("${path.module}", "pkg/controller/*.go")
  loaders_files = fileset("${path.module}", "pkg/loaders/*.go")
  logger_files = fileset("${path.module}", "pkg/logger/*.go")
  prom_files = fileset("${path.module}", "pkg/prom/*.go")
  types_files = fileset("${path.module}", "pkg/types/*.go")
  yace_files = fileset("${path.module}", "pkg/yace/*.go")
  main_files = fileset("${path.module}", "main.go")

  all_files = concat(
    tolist(local.controller_files),
    tolist(local.loaders_files),
    tolist(local.logger_files),
    tolist(local.prom_files),
    tolist(local.types_files),
    tolist(local.yace_files),
    tolist(local.main_files)
  )

  //files_unsorted = fileset("${path.module}/../pkg/**", "*.go")
  files_sorted = sort(local.all_files)
  aggregated_hash = join("", [
    for f in local.files_sorted : md5(file("${path.module}/../${f}"))
  ])
}