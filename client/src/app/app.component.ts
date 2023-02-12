import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: `
    <!--The content below is only a placeholder and can be replaced.-->
    <div style="text-align:center" class="content">
      <span style="display: block">{{ title }} app is running!</span>
    </div>
    <main>
      <router-outlet> </router-outlet>
    </main>
  `,
  styles: [],
})
export class AppComponent {
  title = 'activitypub.lacolaco.net';
}
