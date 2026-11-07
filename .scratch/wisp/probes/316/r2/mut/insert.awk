NR==150 {
  print
  printf "\t// 316-r2 M-shape replay mutation - overlay ONLY, never the working tree:\n"
  printf "\t// the start-up banner appends one operator-sentence-shaped line carrying the\n"
  printf "\t// very words of the restart notice: the exact M shape ticket 316 names.\n"
  printf "\tfmt.Fprintf(rt.stdout,\n"
  printf "\t\t\"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）。原因：上一版横幅（316-r2 M 形复造）。涉及：app.autostart / app.language / app.single_instance。注意两件事都没发生：文件里的新值没有被丢掉，也没有被偷偷生效。\\n\")\n"
  next
}
{ print }
