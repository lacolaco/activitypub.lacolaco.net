import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AdminApiClient } from '../../shared/api';

@Component({
  selector: 'app-create-user',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, FormsModule],
  template: `
    <form (ngSubmit)="createUser()" class="flex flex-col gap-y-2">
      <label class="grid grid-cols-[100px_1fr] items-center justify-stretch gap-x-2">
        <span>host</span>
        <input type="text" name="host" [ngModel]="host()" (ngModelChange)="host.set($event)" />
      </label>
      <label class="grid grid-cols-[100px_1fr] items-center justify-stretch gap-x-2">
        <span>username</span>
        <input type="text" name="username" [ngModel]="username()" (ngModelChange)="username.set($event)" />
      </label>
      <label class="grid grid-cols-[100px_1fr] items-center justify-stretch gap-x-2">
        <span>displayName</span>
        <input type="text" name="displayName" [ngModel]="displayName()" (ngModelChange)="displayName.set($event)" />
      </label>
      <label class="grid grid-cols-[100px_1fr] items-center justify-stretch gap-x-2">
        <span> description</span>
        <input type="text" name="description" [ngModel]="description()" (ngModelChange)="description.set($event)" />
      </label>
      <label class="grid grid-cols-[100px_1fr] items-center justify-stretch gap-x-2">
        <span>url</span>
        <input type="text" name="url" [ngModel]="url()" (ngModelChange)="url.set($event)" />
      </label>

      <button type="submit" class="border border-gray-700 p-2">create user</button>
    </form>
  `,
  host: { class: 'block' },
})
export class CreateUserComponent {
  #api = inject(AdminApiClient);

  host = signal('test.lacolaco.social');
  username = signal('lacolaco');
  displayName = signal('lacolaco');
  description = signal('');
  url = signal('https://lacolaco.net');

  async createUser() {
    const user = await this.#api.createUser({
      host: this.host(),
      username: this.username(),
      displayName: this.displayName(),
      description: this.description(),
      url: this.url(),
      icon: { url: 'https://github.com/lacolaco.png' },
      attachments: [],
    });

    console.log(user);
  }
}
