<script lang="ts">
  import { PILL_SOLID } from "../ui/buttons";
  import { ICON } from "../ui/icons";
  import Icon from "./Icon.svelte";

  // The closing nudge at the bottom of the content pages: back to checking a link.
</script>

<section
  class="mt-24 pt-12 border-t border-gray-200 dark:border-gray-800 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6"
>
  <p class="font-serif text-3xl md:text-4xl tracking-[-0.01em] leading-tight">
    遇到拿不准的<span class="italic">链接？</span>
  </p>
  <a href="/" class="{PILL_SOLID} px-6 py-3 self-start sm:self-auto shrink-0">
    立即检查
    <Icon path={ICON.arrowRight} />
  </a>
</section>
