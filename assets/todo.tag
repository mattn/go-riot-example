
<todo>

  <h3>{ opts.title }</h3>

  <ul>
    <li each={ items.filter(filter) }>
      <label class={ completed: done }>
        <input type="checkbox" checked={ done } onclick={ parent.toggle }> { title }
      </label>
    </li>
  </ul>

  <form onsubmit={ add }>
    <input name="input" onkeyup={ edit }>
    <button disabled={ !text }>Add #{ items.filter(filter).length + 1 }</button>
  </form>

  <!-- this script tag is optional -->
  <script>
    this.items = opts.items

    edit(e) {
      this.text = e.target.value
    }

    add(e) {
      var item = this;
      if (!item.text) return;
      request
        .post("/api")
        .type('form')
        .send({title: item.text})
        .end(function() {
          item.items.push({ title: item.text })
          item.text = item.input.value = ''
          item.update();
        });
    }

    filter(item) {
      return !item.hidden
    }

    toggle(e) {
      var item = e.item
      item.done = !item.done
      request
        .post("/api/" + item.id)
        .type('form')
        .send(item)
        .end(function() {
          riot.update();
        });
      return true
    }
  </script>

</todo>
