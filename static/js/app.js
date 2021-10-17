function notifyModal(icon, title, text, confirmBtnText) {
    Swal.fire({
      icon: icon,
      title: title,
      html: text,
      confirmButtonText: confirmBtnText,
    });
  }
  function Prompt() {
    const toast = (c) => {
      const { msg = '', icon = 'success', position = 'top-end' } = c;

      const Toast = Swal.mixin({
        toast: true,
        icon: icon,
        title: msg,
        position: position,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.addEventListener('mouseenter', Swal.stopTimer);
          toast.addEventListener('mouseleave', Swal.resumeTimer);
        },
      });

      Toast.fire({});
    };
    const success = (c) => {
      const { msg = '', title = '', footer = '' } = c;
      Swal.fire({
        icon: 'success',
        title: title,
        text: msg,
        footer: footer,
      });
    };
    const error = (c) => {
      const { msg = '', title = '', footer = '' } = c;
      Swal.fire({
        icon: 'error',
        title: title,
        text: msg,
        footer: footer,
      });
    };
    const custom = async (c) => {
      const {
        icon = '',
        msg = '',
        title = '',
        showConfirmButton = true,
      } = c;
      const { value: result } = await Swal.fire({
        icon: icon,
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: true,
        showConfirmButton: showConfirmButton,
        didOpen: () => {
          if (c.didOpen !== undefined) {
            c.didOpen();
          }
        },
        // preConfirm: () => {
        //   return [
        //     document.getElementById('start-modal').value,
        //     document.getElementById('end-modal').value,
        //   ];
        // },
      });

      if (result) {
        if (result.dismiss !== Swal.DismissReason.cancel) {
          if (result.value !== '') {
            if (c.callback !== undefined) {
              c.callback(result);
            } else {
              c.callback(false);
            }
          } else {
            c.callback(false);
          }
        }
      }
    };
    return {
      toast,
      success,
      error,
      custom,
    };
  }