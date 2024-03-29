import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { AppStrokedButton } from './shared/ui/button';
import { SearchComponent } from './search/search.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterOutlet, AppStrokedButton, SearchComponent],
  template: `
    <header class="p-4 shadow-gray-500 shadow-sm z-10">
      <div class="container flex flex-row items-center justify-between">
        <h1 class="font-bold">
          <a routerLink="/">console.lacolaco.social</a>
        </h1>
      </div>
    </header>
    <main class="flex-auto container py-4 flex flex-col gap-y-2">
      <app-search-remote-user></app-search-remote-user>
      <router-outlet></router-outlet>
    </main>
  `,
  host: { class: 'flex flex-col w-full h-full bg-white font-sans' },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppComponent {}
