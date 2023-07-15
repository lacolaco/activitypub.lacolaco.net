import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { AdminApiClient } from '../shared/api';
import { ActivityPubPerson } from '../shared/models';
import { AppStrokedButton } from '../shared/ui/button';
import { FormFieldModule } from '../shared/ui/form-field';

@Component({
  selector: 'app-search-remote-user',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormFieldModule, AppStrokedButton],
  template: `
    <div class="flex-auto flex flex-col justify-start gap-y-2">
      <h2 class="text-xl">Search remote user</h2>
      <form [formGroup]="form" (submit)="searchUser()">
        <app-form-field label="@username@hostname" [showLabel]="true">
          <input type="text" placeholder="@username@hostname" formControlName="userId" />
        </app-form-field>
      </form>

      <div *ngIf="searched()">
        <div *ngIf="person() as p" class="flex flex-col items-start rounded-lg bg-panel p-4 shadow">
          <div *ngIf="p.icon">
            <img [src]="p.icon.url" class="w-24 h-24 rounded-lg" />
          </div>
          <div class="flex flex-col items-start py-2">
            <span class="font-bold text-xl">{{ p.name }}</span>
            <span class="text-xs break-all text-gray-600">{{ p.id }}</span>
          </div>
          <div class="py-2" [innerHTML]="p.summary"></div>
        </div>

        <div *ngIf="!person()">
          <p>Not found</p>
        </div>
      </div>
    </div>
  `,
  host: { class: 'block' },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent {
  private readonly api = inject(AdminApiClient);

  readonly person = signal<ActivityPubPerson | null>(null);
  readonly searched = signal(false);

  readonly form = new FormGroup({
    userId: new FormControl('@lacolaco@misskey.io', { nonNullable: true }),
  });

  async searchUser() {
    const userId = this.form.getRawValue().userId;
    if (!userId) {
      return;
    }
    try {
      const person = await this.api.searchRemotePerson(userId);
      this.person.set(person);
      this.searched.set(true);
      console.log(person);
    } catch (e) {
      console.error(e);
    }
  }
}
