import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { firstValueFrom, lastValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';
import { AppStrokedButton } from '../shared/ui/button';
import { FormFieldModule } from '../shared/ui/form-field';

type ActivityPubPerson = {
  id: string;
  name: string;
  inbox: string;
  summary: string;
  icon?: { url: string };
};

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

          <div class="flex flex-row gap-x-2">
            <button app-stroked-button class="bg-white" (click)="requestFollow(p)">Follow</button>
            <button app-stroked-button class="bg-white" (click)="requestUnfollow(p)">Unfollow</button>
          </div>
        </div>

        <div *ngIf="!person()">
          <p>Not found</p>
        </div>
      </div>
    </div>
  `,
  host: { class: 'flex flex-col gap-y-2 h-full' },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent {
  private readonly http = inject(HttpClient);

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
      const person = await firstValueFrom(
        this.http.get<ActivityPubPerson>(`${environment.backend}/admin/search/person/${userId}`),
      );
      this.person.set(person);
      this.searched.set(true);
      console.log(person);
    } catch (e) {
      console.error(e);
    }
  }

  async requestFollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`${environment.backend}/admin/following/create`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }

  async requestUnfollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`${environment.backend}/admin/following/delete`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }
}
