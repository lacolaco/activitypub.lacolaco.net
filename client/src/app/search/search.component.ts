import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <h1>Search</h1>

    <form>
      <input type="text" placeholder="@username@hostname" />
    </form>
  `,
  styles: [],
})
export class SearchComponent {}
