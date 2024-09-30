" Use Vim-Plug to manage plugins
call plug#begin('~/.vim/plugged')

" Python autocompletion
Plug 'davidhalter/jedi-vim'

" File explorer
Plug 'preservim/nerdtree'

" Git integration
Plug 'airblade/vim-gitgutter'

" Status/tab line
Plug 'vim-airline/vim-airline'

" End of plugins
call plug#end()

" NERDTree settings
map <C-n> :NERDTreeToggle<CR>

" Jedi settings
let g:jedi#completions_enabled = 1
let g:jedi#use_tabs_not_buffers = 1

" GitGutter settings
let g:gitgutter_enabled = 1

" Airline settings
let g:airline#extensions#tabline#enabled = 1

