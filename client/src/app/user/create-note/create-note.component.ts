import { ChangeDetectionStrategy, Component, Input, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LocalUser } from 'src/app/shared/models';
import { FormsModule } from '@angular/forms';
import { AdminApiClient } from '../../shared/api';

@Component({
  selector: 'app-create-note',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './create-note.component.html',
  host: { class: 'block' },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CreateNoteComponent {
  api = inject(AdminApiClient);

  @Input() user!: LocalUser;

  content = signal<string>('');

  async submit() {
    console.log('submit');
    const note = {
      content: this.content(),
    };

    await this.api.postUserNote(this.user, note);
  }
}
