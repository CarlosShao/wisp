/* ============================================================================
   Config base - the fingerprint half of the config editor's drift handshake
   (ticket 39 AC#4), in one place.
   ----------------------------------------------------------------------------
   THE MECHANISM (optimistic concurrency, the standard If-Match shape): the
   config.get read face carries `baseSHA` - the disk file's sha256 at read
   time (internal/config/fingerprint_39.go FileSHA256, through the payload's
   baseSHA key). Every config.set this GUI sends echoes that fingerprint back
   on the envelope's configBaseSHA key; the write leg compares it against the
   file's OWN hash right now and refuses with refused-drift (盘上已变) on a
   mismatch, so a stale form cannot silently replace the external edit it
   never saw.

   WHY A MODULE CACHE AND NOT PROP PLUMBING: the writes leave from several
   pages (the appearance rows, the model page, the composer switches, the
   plugin inventory) and the reads land in several hooks (useConfigView).
   Each of those already re-reads the whole view after every accepted write
   (the reload event, PLAN.md:1044's no-local-copy rule) - the cache here is
   NOT a second copy of the config VALUES, it is one fact about the DISK that
   the page saw, exactly what an If-Match token is. It advances on:
     - every successful config.get parse (config-view.ts useConfigView), and
     - every accepted write's reload event (config-actions.ts), which carries
       the post-write fingerprint (baseSHA) so back-to-back saves never trip
       the drift check against their own predecessor.
   It never advances on a refusal: the refusal IS the stale view.

   ABSENT = NO CHECK. undefined (no read yet, or a leg that did not answer
   the fingerprint half) means the envelope carries no configBaseSHA key -
   byte-identically to the pre-39 shape. A page that never received a
   fingerprint cannot be accused of holding a stale one, and the write runs
   without the check rather than refusing a page that was never told.
   ============================================================================ */

let lastBaseSHA: string | undefined;

/** Record the fingerprint one successful read (or one accepted write's reload
    event) answered. undefined clears it - "no fingerprint held". */
export function rememberConfigBaseSHA(sha: string | undefined): void {
  lastBaseSHA = typeof sha === "string" && sha.trim() !== "" ? sha : undefined;
}

/** The fingerprint the GUI's freshest view of the disk answered, or undefined
    when no read has answered one yet (the envelope then carries no
    configBaseSHA key at all). */
export function currentConfigBaseSHA(): string | undefined {
  return lastBaseSHA;
}
