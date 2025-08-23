import Swal from 'sweetalert2';

export async function showAlert(fn,opts = {}) {
  return await Swal.fire({
        title: opts['title'] || 'Are you sure?',
        text: opts['text'] || "You won't be able to undo this!",
        icon: opts['icon'] || 'warning',
        showCancelButton: true,
        confirmButtonText: opts['confirmButtonText'] || 'Yes',
        cancelButtonText: 'Cancel',
        customClass: {
          confirmButton: 'btn btn-primary btn-lg mr-2',
          cancelButton: 'btn btn-danger btn-lg',
          loader: 'custom-loader',
        },
        loaderHtml: '<div class="spinner-border text-primary"></div>',
        preConfirm: async () => {
          Swal.showLoading();
          try {
            await fn()
          } catch (error) {
            Swal.showValidationMessage(`Request failed: ${error}`);
          }
        },
        allowOutsideClick: () => !Swal.isLoading() // prevent closing while loading
      })
}