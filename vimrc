" Example vim setup to use edit a compacted format but save the normal format
function! GoFormatFile()
    if executable("gofmt")
        !gofmt -w <afile>
    endif
endfunction

function! GoCompactBuffer()
    let s:fmt = "gocompact"
    call s:formatBuffer()
endfunction

function! s:findFormatter()
    return executable(s:fmt)
endfunction

function! s:formatBuffer()
    if s:findFormatter() && (&modifiable == 1)
        let l:curw=winsaveview()
        let l:tmpname=tempname()
        call writefile(getline(1,'$'), l:tmpname)
        call system(s:fmt . " " . l:tmpname ." > /dev/null 2>&1")
        if v:shell_error == 0
            try | silent undojoin | catch | endtry
            silent exe '%!' . s:fmt
        endif
        call delete(l:tmpname)
        call winrestview(l:curw)
    endif
endfunction

augroup gofmtbuffer
    au BufUnload                *.go silent call GoFormatFile()
    au BufReadPost,BufWrite     *.go silent call GoCompactBuffer()
augroup END
