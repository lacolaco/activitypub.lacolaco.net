import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, effect, inject, Input, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';

export type LocalUser = {
  id: string;
  name: string;
  description: string;
  icon: { url: string };
  attachments: Array<{ name: string; value: string }>;
};

@Component({
  selector: 'app-user',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div *ngIf="user() as u" class="flex flex-col items-start rounded-lg bg-panel p-4 shadow">
      <div>
        <img [src]="u.icon.url" class="w-24 h-24 rounded-lg" />
      </div>
      <div class="flex flex-col items-start py-2">
        <span class="font-bold text-xl">{{ u.name }}</span>
        <span class="text-sm text-gray-600">@{{ u.id }}@{{ hostname }}</span>
      </div>
      <div class="py-2" [innerHTML]="u.description"></div>
      <div class="w-full">
        <table class="w-full">
          <tr *ngFor="let attachment of u.attachments">
            <td class="font-bold">{{ attachment.name }}</td>
            <td class="ml-2" [innerHTML]="attachment.value"></td>
          </tr>
        </table>
      </div>
    </div>
  `,
  styles: [],
})
export class UserComponent {
  private readonly http = inject(HttpClient);
  readonly user = signal<LocalUser | null>(null);

  readonly hostname = window.location.hostname;

  readonly #username = signal<string | null>(null);

  @Input() set username(username: string) {
    this.#username.set(username);
  }

  constructor() {
    effect(async () => {
      const username = this.#username();
      if (username) {
        const resp = await firstValueFrom(
          this.http.get<{ user: LocalUser }>(`${environment.backend}/api/users/show/${username}`),
        );
        this.user.set(resp.user);
      }
    });
  }
}
