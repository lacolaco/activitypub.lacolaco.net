import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { SearchComponent } from '../search/search.component';
import { AdminApiClient } from '../shared/api';
import { LocalUser } from '../shared/models';
import { CreateUserComponent } from './create-user/create-user.component';

@Component({
  selector: 'app-users',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink, SearchComponent, CreateUserComponent],
  template: `
    <div class="flex flex-col gap-y-2">
      <h2 class="text-xl">Users</h2>
      <ul class="flex flex-col gap-y-2">
        <li *ngFor="let user of users()" class="p-4 rounded border border-gray-300">
          <div>
            <a [routerLink]="['/users', user.host, user.username]">@{{ user.username }}/{{ user.id }}</a>
          </div>
        </li>
      </ul>

      <div>
        <h3 class="text-xl">create new user</h3>

        <app-create-user></app-create-user>
      </div>
    </div>
  `,
  host: { class: 'block' },
})
export class UsersComponent {
  private readonly api = inject(AdminApiClient);

  readonly users = signal<LocalUser[]>([]);

  async ngOnInit() {
    const users = await this.api.getUsers();
    this.users.set(users);
  }
}
