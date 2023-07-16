import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject, Input, signal } from '@angular/core';
import { SearchComponent } from '../search/search.component';
import { AdminApiClient } from '../shared/api';
import { LocalUser } from '../shared/models';

@Component({
  selector: 'app-user',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, SearchComponent],
  template: `
    <div *ngIf="user() as u" class="flex flex-col items-start gap-y-2">
      <h2 class="text-lg">@{{ u.username }}@{{ u.host }}</h2>
      <details class="w-full rounded-lg bg-panel p-4 shadow">
        <summary class="text-md">Raw JSON</summary>
        <pre class="w-full font-mono text-sm overflow-auto">{{ userJSON() }}</pre>
      </details>
      <div>
        <h3>create new note</h3>
      </div>
    </div>
  `,
  styles: [],
})
export class UserComponent {
  private readonly api = inject(AdminApiClient);
  readonly user = signal<LocalUser | null>(null);

  readonly userJSON = computed(() => JSON.stringify(this.user(), null, 2));

  readonly #hostname = signal<string | null>(null);

  @Input() set hostname(hostname: string) {
    this.#hostname.set(hostname);
  }

  readonly #username = signal<string | null>(null);

  @Input() set username(username: string) {
    this.#username.set(username);
  }

  constructor() {
    effect(async () => {
      const hostname = this.#hostname();
      const username = this.#username();
      if (hostname && username) {
        const user = await this.api.getUserByUsername(hostname, username);
        this.user.set(user);
      }
    });
  }
}
