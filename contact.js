function tipeData() {
  let tipeName = document.getElementById("input-name").value;
  let tipeEmail = document.getElementById("input-email").value;
  let tipePhone = document.getElementById("input-phone").value;
  let tipeSubject = document.getElementById("input-subject").value;
  let tipeMessage = document.getElementById("input-message").value;

  console.log(tipeName);
  console.log(tipeEmail);
  console.log(tipePhone);
  console.log(tipeSubject);
  console.log(tipeMessage);

  if (tipeName == "") {
    return alert("Nama wajib diisi");
  }

  if (tipeEmail == "") {
    return alert("Email wajib diisi");
  }

  if (tipePhone == "") {
    return alert("Phone wajib diisi");
  }

  if (tipeSubject == "") {
    return alert("Subject wajib diisi");
  }

  if (tipeMessage == "") {
    return alert("Message wajib diisi");
  }

  let emailReceiver = "ahsanu030721@gmail.com";

  let a = document.createElement("a");

  a.href = `mailto:${emailReceiver}?subject= ${tipeSubject}&body= halo nama saya ${tipeName} no hp saya ${tipePhone} ${tipeMessage}`;

  a.click();
}
