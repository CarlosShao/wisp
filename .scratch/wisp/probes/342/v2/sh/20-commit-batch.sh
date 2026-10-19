#!/bin/sh
# 342-v2 自己的提交尺件（★逐枚字面 pathspec、⛔ 从共享 index 反推、⛔ add -A、⛔ commit -- .）
# 用法：sh 20-commit-batch.sh <批次号 2|3|4>
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
BATCH="$1"
case "$BATCH" in
2)
  git add \
    docs/evidence/s1/342-v2-verdict.md \
    .scratch/wisp/probes/342/v2/output-budget.md \
    .scratch/wisp/probes/342/v2/readings-windows-qualify.md \
    .scratch/wisp/probes/342/v2/logs/four-counts-probe.tsv \
    .scratch/wisp/probes/342/v2/logs/four-counts-windows-fullpkg.tsv \
    .scratch/wisp/probes/342/v2/logs/gate-d22scan.txt \
    .scratch/wisp/probes/342/v2/logs/gate-length.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-darwin.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-linux.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-local.txt \
    .scratch/wisp/probes/342/v2/logs/ii02-mut-W-treeB-probe.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii02-mut-W-treeL-probe.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii03-mut-X-open-treeL-1.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii03-mut-X-open-treeL-2.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-A-unlink-byname.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-B-refuse-every-name.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-B-vet.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeB-winfullpkg-1.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeB-winfullpkg-2.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeL-winfullpkg-1.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeL-winfullpkg-2.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/logs-12-windows-fullpkg.txt \
    .scratch/wisp/probes/342/v2/logs/probe-342v2-windows_test.go.txt \
    .scratch/wisp/probes/342/v2/logs/probe-err-unix_test.go.txt \
    .scratch/wisp/probes/342/v2/logs/probe2-run.txt \
    .scratch/wisp/probes/342/v2/logs/probe3-run.txt \
    .scratch/wisp/probes/342/v2/logs/probe4-run.txt \
    .scratch/wisp/probes/342/v2/logs/v02-errclass-treeB.raw.txt \
    .scratch/wisp/probes/342/v2/logs/v02-errclass-treeL.raw.txt \
    .scratch/wisp/probes/342/v2/logs/v03-classify-decision-points.txt \
    .scratch/wisp/probes/342/v2/logs/x01-discipline.txt \
    .scratch/wisp/probes/342/v2/sh/10-teeth-A-unlink-byname.sh \
    .scratch/wisp/probes/342/v2/sh/11-teeth-B-refuse-every-name.sh \
    .scratch/wisp/probes/342/v2/sh/12-windows-fullpkg-pair.sh \
    .scratch/wisp/probes/342/v2/sh/13-windows-probe.sh \
    .scratch/wisp/probes/342/v2/sh/14-unplug-production.sh \
    .scratch/wisp/probes/342/v2/sh/15-errclassify-pair.sh \
    .scratch/wisp/probes/342/v2/sh/16-classify-decision-points.sh \
    .scratch/wisp/probes/342/v2/sh/17-discipline.sh \
    .scratch/wisp/probes/342/v2/sh/18-output-budget.sh \
    .scratch/wisp/probes/342/v2/sh/19-verdict-cell-lines.sh
  git commit -F 'D:\tmp\wisp\20261019\342v2\cm-02.txt' -- \
    docs/evidence/s1/342-v2-verdict.md \
    .scratch/wisp/probes/342/v2/output-budget.md \
    .scratch/wisp/probes/342/v2/readings-windows-qualify.md \
    .scratch/wisp/probes/342/v2/logs/four-counts-probe.tsv \
    .scratch/wisp/probes/342/v2/logs/four-counts-windows-fullpkg.tsv \
    .scratch/wisp/probes/342/v2/logs/gate-d22scan.txt \
    .scratch/wisp/probes/342/v2/logs/gate-length.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-darwin.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-linux.txt \
    .scratch/wisp/probes/342/v2/logs/gate-vet-local.txt \
    .scratch/wisp/probes/342/v2/logs/ii02-mut-W-treeB-probe.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii02-mut-W-treeL-probe.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii03-mut-X-open-treeL-1.raw.txt \
    .scratch/wisp/probes/342/v2/logs/ii03-mut-X-open-treeL-2.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-A-unlink-byname.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-B-refuse-every-name.raw.txt \
    .scratch/wisp/probes/342/v2/logs/iii-B-vet.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeB-winfullpkg-1.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeB-winfullpkg-2.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeL-winfullpkg-1.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/iv02-treeL-winfullpkg-2.rednames.txt \
    .scratch/wisp/probes/342/v2/logs/logs-12-windows-fullpkg.txt \
    .scratch/wisp/probes/342/v2/logs/probe-342v2-windows_test.go.txt \
    .scratch/wisp/probes/342/v2/logs/probe-err-unix_test.go.txt \
    .scratch/wisp/probes/342/v2/logs/probe2-run.txt \
    .scratch/wisp/probes/342/v2/logs/probe3-run.txt \
    .scratch/wisp/probes/342/v2/logs/probe4-run.txt \
    .scratch/wisp/probes/342/v2/logs/v02-errclass-treeB.raw.txt \
    .scratch/wisp/probes/342/v2/logs/v02-errclass-treeL.raw.txt \
    .scratch/wisp/probes/342/v2/logs/v03-classify-decision-points.txt \
    .scratch/wisp/probes/342/v2/logs/x01-discipline.txt \
    .scratch/wisp/probes/342/v2/sh/10-teeth-A-unlink-byname.sh \
    .scratch/wisp/probes/342/v2/sh/11-teeth-B-refuse-every-name.sh \
    .scratch/wisp/probes/342/v2/sh/12-windows-fullpkg-pair.sh \
    .scratch/wisp/probes/342/v2/sh/13-windows-probe.sh \
    .scratch/wisp/probes/342/v2/sh/14-unplug-production.sh \
    .scratch/wisp/probes/342/v2/sh/15-errclassify-pair.sh \
    .scratch/wisp/probes/342/v2/sh/16-classify-decision-points.sh \
    .scratch/wisp/probes/342/v2/sh/17-discipline.sh \
    .scratch/wisp/probes/342/v2/sh/18-output-budget.sh \
    .scratch/wisp/probes/342/v2/sh/19-verdict-cell-lines.sh
  ;;
3)
  git add \
    docs/evidence/s1/348-v1-verdict.md \
    .scratch/wisp/probes/348/v1/readings-deleted-leg.md \
    .scratch/wisp/probes/348/v1/logs/ii02-mut-W-treeB-probe.raw.txt \
    .scratch/wisp/probes/348/v1/logs/ii02-mut-W-treeL-probe.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-A-unlink-byname.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-B-refuse-every-name.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-B-vet.txt \
    .scratch/wisp/probes/348/v1/logs/probe-steady-state-windows_test.go.txt \
    .scratch/wisp/probes/348/v1/logs/readings-posix-nails.md \
    .scratch/wisp/probes/348/v1/logs/readings-windows-qualify.md \
    .scratch/wisp/probes/348/v1/sh/01-roster-two-anchors.sh \
    .scratch/wisp/probes/348/v1/sh/02-leg-diffnamestatus.sh \
    .scratch/wisp/probes/348/v1/sh/03-container-run.sh \
    .scratch/wisp/probes/348/v1/sh/04-build-trees.sh \
    .scratch/wisp/probes/348/v1/sh/05-container-run.sh \
    .scratch/wisp/probes/348/v1/sh/06-posix-nails.sh \
    .scratch/wisp/probes/348/v1/sh/07-errnotexist-consumers.sh \
    .scratch/wisp/probes/348/v1/sh/08-windows-nails.sh \
    .scratch/wisp/probes/348/v1/sh/09-posix-fullpkg-pair.sh \
    .scratch/wisp/probes/348/v1/sh/10-teeth-A-unlink-byname.sh \
    .scratch/wisp/probes/348/v1/sh/11-teeth-B-refuse-every-name.sh \
    .scratch/wisp/probes/348/v1/sh/12-windows-fullpkg-pair.sh \
    .scratch/wisp/probes/348/v1/sh/13-windows-probe.sh \
    .scratch/wisp/probes/348/v1/sh/14-unplug-production.sh \
    .scratch/wisp/probes/348/v1/sh/15-errclassify-pair.sh \
    .scratch/wisp/probes/348/v1/sh/16-classify-decision-points.sh \
    .scratch/wisp/probes/348/v1/sh/17-discipline.sh \
    .scratch/wisp/probes/348/v1/sh/18-output-budget.sh \
    .scratch/wisp/probes/348/v1/sh/19-verdict-cell-lines.sh \
    .scratch/wisp/probes/348/v1/sh/20-commit-batch.sh
  git commit -F 'D:\tmp\wisp\20261019\342v2\cm-03.txt' -- \
    docs/evidence/s1/348-v1-verdict.md \
    .scratch/wisp/probes/348/v1/readings-deleted-leg.md \
    .scratch/wisp/probes/348/v1/logs/ii02-mut-W-treeB-probe.raw.txt \
    .scratch/wisp/probes/348/v1/logs/ii02-mut-W-treeL-probe.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-A-unlink-byname.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-B-refuse-every-name.raw.txt \
    .scratch/wisp/probes/348/v1/logs/iii-B-vet.txt \
    .scratch/wisp/probes/348/v1/logs/probe-steady-state-windows_test.go.txt \
    .scratch/wisp/probes/348/v1/logs/readings-posix-nails.md \
    .scratch/wisp/probes/348/v1/logs/readings-windows-qualify.md \
    .scratch/wisp/probes/348/v1/sh/01-roster-two-anchors.sh \
    .scratch/wisp/probes/348/v1/sh/02-leg-diffnamestatus.sh \
    .scratch/wisp/probes/348/v1/sh/03-container-run.sh \
    .scratch/wisp/probes/348/v1/sh/04-build-trees.sh \
    .scratch/wisp/probes/348/v1/sh/05-container-run.sh \
    .scratch/wisp/probes/348/v1/sh/06-posix-nails.sh \
    .scratch/wisp/probes/348/v1/sh/07-errnotexist-consumers.sh \
    .scratch/wisp/probes/348/v1/sh/08-windows-nails.sh \
    .scratch/wisp/probes/348/v1/sh/09-posix-fullpkg-pair.sh \
    .scratch/wisp/probes/348/v1/sh/10-teeth-A-unlink-byname.sh \
    .scratch/wisp/probes/348/v1/sh/11-teeth-B-refuse-every-name.sh \
    .scratch/wisp/probes/348/v1/sh/12-windows-fullpkg-pair.sh \
    .scratch/wisp/probes/348/v1/sh/13-windows-probe.sh \
    .scratch/wisp/probes/348/v1/sh/14-unplug-production.sh \
    .scratch/wisp/probes/348/v1/sh/15-errclassify-pair.sh \
    .scratch/wisp/probes/348/v1/sh/16-classify-decision-points.sh \
    .scratch/wisp/probes/348/v1/sh/17-discipline.sh \
    .scratch/wisp/probes/348/v1/sh/18-output-budget.sh \
    .scratch/wisp/probes/348/v1/sh/19-verdict-cell-lines.sh \
    .scratch/wisp/probes/348/v1/sh/20-commit-batch.sh
  ;;
esac
echo "COMMIT_BATCH_RC=$?"
