import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject, Input, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AdminApiClient } from '../shared/api';
import { LocalUser } from '../shared/models';
import { CreateNoteComponent } from './create-note/create-note.component';

@Component({
  selector: 'app-user',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink, CreateNoteComponent],
  template: `
    <div class="flex flex-col items-start gap-y-2">
      <header class="flex flex-row gap-x-1 text-xl">
        <a routerLink="/users">Users</a>
        <span class="text-gray-500">/</span>
        <h2>@{{ username }}@{{ hostname }}</h2>
      </header>

      <ng-container *ngIf="user() as u">
        <details class="w-full rounded-lg bg-panel p-4 shadow">
          <summary class="text-sm">Raw JSON</summary>
          <pre class="w-full font-mono text-xs overflow-auto">{{ userJSON() }}</pre>
        </details>

        <div>
          <h3>Edit user data</h3>
        </div>

        <div class="flex flex-col w-full">
          <h3 class="font-bold text-lg">Create new note</h3>
          <app-create-note [user]="u" class="w-full"></app-create-note>
        </div>
      </ng-container>
    </div>
  `,
  styles: [],
})
export class UserComponent {
  private readonly api = inject(AdminApiClient);
  readonly user = signal<LocalUser | null>(null);

  readonly userJSON = computed(() => JSON.stringify(this.user(), null, 2));

  readonly #hostname = signal<string | null>(null);

  @Input() set hostname(hostname: string | null) {
    this.#hostname.set(hostname);
  }
  get hostname() {
    return this.#hostname();
  }

  readonly #username = signal<string | null>(null);

  @Input() set username(username: string | null) {
    this.#username.set(username);
  }
  get username() {
    return this.#username();
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
