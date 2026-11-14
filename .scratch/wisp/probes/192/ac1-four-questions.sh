#!/bin/sh
# Ticket 192 AC#1 -- the four questions, each answered by a code position and
# by the criterion that pins it. This script PRINTS the positions (grep -n on
# the file that ships, no paraphrase) and RUNS the criteria, so the answers in
# the evidence file cannot drift from the code without this log changing.
#
#   Q1 which classes are previewable      -> classifyPreviewArtifact + the four kinds
#   Q2 path or bytes                      -> the guard, the artifacts dir, "no path in the response"
#   Q3 truncation semantics for big files -> textFirstView/textWindow + the class caps
#   Q4 what the user sees when truncated  -> the four notice lines the first view builds
#
# Usage: sh .scratch/wisp/probes/192/ac1-four-questions.sh
set -u
here=$(cd "$(dirname "$0")" && pwd)
repo=$(cd "$here/../../../.." && pwd)
logs="$here/logs"
mkdir -p "$logs"
cd "$repo" || exit 2
f=internal/panel/preview_192.go
t=internal/panel/preview_192_test.go
log="$logs/11-ac1-four-questions.log"

{
	echo "# AC#1 four questions -- code positions and the criteria that pin them"
	echo "# file: $f (production) / $t (criteria)"
	echo
	echo "## Q1 能预览哪几类"
	echo '$ grep -n "PreviewKind" '"$f"
	grep -n 'PreviewKind\(HTML\|Image\|PDF\|Text\) \+PreviewKind' "$f"
	echo '$ grep -n "func classifyPreviewArtifact" '"$f"
	grep -n 'func classifyPreviewArtifact' "$f"
	echo '$ grep -n "attachmentTypes" '"$f"
	grep -n 'range attachmentTypes' "$f"
	echo '# criterion: TestPreview192FourClassesAndTheRestRefused / TestPreview192ClassComesFromBytesNotFromTheName'
	echo
	echo "## Q2 渲染本地路径还是内容字节"
	echo '$ grep -n "r.guard(name)" '"$f"
	grep -n 'r.guard(name)' "$f"
	echo '$ grep -n "filepath.Join" '"$f"
	grep -n 'filepath.Join' "$f"
	echo '$ grep -n "ArtifactsDir" '"$f"
	grep -n 'ArtifactsDir()' "$f"
	echo '# criterion: TestPreview192RefusesPathShapedNamesBeforeItOpensAnything / TestPreview192TheResponseNeverCarriesTheArtifactsPath'
	echo
	echo "## Q3 大文件截断语义"
	echo '$ grep -n "previewHeadBytes\|previewTailBytes\|PreviewMaxTextWindow\|preview.*Cap" '"$f"
	grep -n 'previewHeadBytes\s*=\|previewTailBytes\s*=\|PreviewMaxTextWindow\s*=\|previewHTMLCap\s*=\|previewImageCap\s*=\|previewPDFCap\s*=' "$f"
	echo '$ grep -n "func (r \*PreviewResolver) text" '"$f"
	grep -n 'func (r \*PreviewResolver) text' "$f"
	echo '# criterion: TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs / TestPreview192ThePointerWalksEveryByteOfTheFile / TestPreview192BinaryOverTheCapIsRefusedNotTruncated / TestPreview192SmallTextComesBackWholeWithNoEnvelope'
	echo
	echo "## Q4 截断时界面看到什么"
	echo '$ grep -n "previewNotice" '"$f"
	grep -n 'previewNotice\(Total\|Resume\|Body\|Elided\) \+=' "$f"
	echo '$ grep -n "Fprintf(&b" '"$f"
	grep -n 'Fprintf(&b' "$f"
	echo '# criterion: TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs (it PARSES the display and checks the numbers add up)'
	echo
	echo "## the criteria, run"
	(cd "$repo" && go test -count=1 -v -run '^(TestPreview192FourClassesAndTheRestRefused|TestPreview192ClassComesFromBytesNotFromTheName|TestPreview192RefusesPathShapedNamesBeforeItOpensAnything|TestPreview192TheResponseNeverCarriesTheArtifactsPath|TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs|TestPreview192ThePointerWalksEveryByteOfTheFile|TestPreview192BinaryOverTheCapIsRefusedNotTruncated|TestPreview192SmallTextComesBackWholeWithNoEnvelope)$' ./internal/panel/ 2>&1) |
		grep -E '^(=== RUN|--- (PASS|FAIL|SKIP)|ok |FAIL|PASS)' | sed 's/^/   /'
	echo
	echo "## AC#2 criteria, run (the injection face)"
	(cd "$repo" && go test -count=1 -v -run '^(TestPreview192HTMLIsSandboxedSoAPreviewCannotSpeak|TestPreview192PreviewIsNotAMessageRoute|TestPreview192QuoteDoorStampsAndTheEgressGateThenCatches|TestPreview192QuoteDoorRefusesTextThatIsNotInTheArtifact|TestPreview192QuoteDoorRefusesWhatItCannotStamp|TestPreview192SourceToolIsTheContractName)$' ./internal/panel/ 2>&1) |
		grep -E '^(=== RUN|--- (PASS|FAIL|SKIP)|ok |FAIL|PASS)' | sed 's/^/   /'
} 2>&1 | tee "$log"
