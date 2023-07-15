import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { ChangeDetectionStrategy, Component, computed, effect, inject, Input, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';
import { SearchComponent } from '../search/search.component';

export type LocalUser = {
  id: string;
  username: string;
  name: string;
  description: string;
  icon: { url: string };
  attachments: Array<{ name: string; value: string }>;
};

@Component({
  selector: 'app-user',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, SearchComponent],
  template: `
    <div *ngIf="user() as u" class="flex flex-col items-start gap-y-2">
      <details class="w-full rounded-lg bg-panel p-4 shadow">
        <summary class="text-md">@{{ u.username }}@{{ hostname }}</summary>
        <pre class="w-full font-mono text-sm overflow-auto">{{ userJSON() }}</pre>
      </details>

      <app-search-remote-user class="w-full"></app-search-remote-user>
    </div>
  `,
  styles: [],
})
export class UserComponent {
  private readonly http = inject(HttpClient);
  readonly user = signal<LocalUser | null>(null);

  readonly userJSON = computed(() => JSON.stringify(this.user(), null, 2));

  readonly hostname = window.location.hostname;

  readonly #username = signal<string | null>(null);

  @Input() set username(username: string) {
    this.#username.set(username);
  }

  constructor() {
    effect(async () => {
      const username = this.#username();
      if (username) {
        const user = await firstValueFrom(
          this.http.get<LocalUser>(`${environment.backend}/admin/users/show/${username}`),
        );
        this.user.set(user);
      }
    });
  }
}
