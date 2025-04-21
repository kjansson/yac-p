locals {
  config_file_contents = file(var.config_file_local_path)
  files_unsorted = fileset("${path.module}/..", "*.go")
  files_sorted = sort(local.files_unsorted)
  aggregated_hash = join("", [
    for f in local.files_sorted : md5(file("${path.module}/../${f}"))
  ])
}