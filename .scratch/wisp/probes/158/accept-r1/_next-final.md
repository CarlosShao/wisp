`next=` **不写死号、也不写死枚数**（本程一度在这两处各栽一次，见 §11.3 A-5）：接手时现量这四件——
① `git rev-parse --short c7f638c` 与被验面是否仍未动：`git diff --name-only c7f638c..HEAD -- '*.go'`（本程取版时＝空）；
② 本程自己的提交枚数与清单：`git log --oneline --grep='票 158 验收 r1'`（**别引用本表任何一枚号**）；
③ 本程提交里有没有混进别人的路径：`git log --no-walk --pretty=tformat: --name-only <那一列号> | sort -u`
   ——**必须按号集合算，不能拿 `c7f638c..HEAD` 区间当尺**（本表 §1.2/§7 记过：区间尺会把编排者的 ledger 写成我的）；
④ 门禁是否还红在同一个字形上：`sh scripts/d22scan.sh; echo rc=$?`（本程两次 rc=1，唯一 finding＝`right-rail.tsx:108`）。
