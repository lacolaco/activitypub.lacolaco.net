import { CommonModule } from '@angular/common';
import { Component, effect, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { Auth, GoogleAuthProvider, authState, signInWithPopup, signOut } from '@angular/fire/auth';
import { RouterLink, RouterOutlet } from '@angular/router';
import { AppStrokedButton } from './shared/ui/button';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterOutlet, AppStrokedButton],
  template: `
    <header class="p-4 shadow-gray-500 shadow-sm z-10">
      <div class="container flex flex-row items-center justify-between">
        <h1 class="font-bold">activitypub.lacolaco.net</h1>
        <div *ngIf="!user()">
          <button app-stroked-button (click)="signIn()" class="text-sm">Sign in</button>
        </div>
        <div *ngIf="user()">
          <button app-stroked-button (click)="signOut()" class="text-sm">Sign out</button>
        </div>
      </div>
    </header>
    <main class="flex-auto container py-4 flex flex-col gap-y-2">
      <div *ngIf="user()" class="flex flex-col">
        <a routerLink="/search" app-stroked-button>Control Room</a>
      </div>
      <router-outlet></router-outlet>
    </main>
  `,
  host: { class: 'flex flex-col w-full h-full bg-white font-sans' },
})
export class AppComponent {
  private readonly auth = inject(Auth);

  readonly user = toSignal(authState(this.auth), { initialValue: null });

  constructor() {
    effect(() => {
      console.log(this.user());
    });
  }

  async signIn() {
    try {
      await signInWithPopup(this.auth, new GoogleAuthProvider());
    } catch (e) {
      console.error(e);
    }
  }

  async signOut() {
    try {
      await signOut(this.auth);
    } catch (e) {
      console.error(e);
    }
  }
}
