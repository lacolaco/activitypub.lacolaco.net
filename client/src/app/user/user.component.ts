import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs';

@Component({
  selector: 'app-user',
  standalone: true,
  imports: [CommonModule],
  template: ` <p>Username: {{ username$ | async }}</p> `,
  styles: [],
})
export class UserComponent {
  private readonly route = inject(ActivatedRoute);

  readonly username$ = this.route.params.pipe(map((params) => params['username']));
}
